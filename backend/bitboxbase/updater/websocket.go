// Copyright 2018 Shift Devices AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package updater

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"

	"github.com/digitalbitbox/bitbox-wallet-app/util/errp"
	"github.com/digitalbitbox/bitbox-wallet-app/util/observable"
	"github.com/digitalbitbox/bitbox-wallet-app/util/observable/action"
	"github.com/flynn/noise"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// runWebsocket sets up loops for sending/receiving, abstracting away the low level details about
// pings, timeouts, clientection closing, etc.
// It returns three channels: one to send messages to the client, one which notifies
// when the clientection was closed and one to receive messages from the client
//
// Closing msg makes runWebsocket's goroutines quit.
// The goroutines close client upon exit, due to a send/receive error or when msg is closed.
// runWebsocket never closes msg.
func (updater *Updater) runWebsocket(client *websocket.Conn) (send chan<- []byte, weHaveQuit chan<- struct{}, receive <-chan []byte, remoteHasQuit <-chan struct{}) {
	const maxMessageSize = 512

	weHaveQuitChan := make(chan struct{})
	remoteHasQuitChan := make(chan struct{})
	sendChan := make(chan []byte)
	receiveChan := make(chan []byte)

	readLoop := func() {
		defer func() {
			close(remoteHasQuitChan)
			_ = client.Close()
		}()
		client.SetReadLimit(maxMessageSize)
		for {
			_, msg, err := client.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					updater.log.WithFields(logrus.Fields{"group": "updater websocket", "error": err}).Error(err.Error())
				}
				break
			}
			messageDecrypted, err := updater.receiveCipher.Decrypt(nil, nil, msg)
			if err != nil {
				updater.log.WithFields(logrus.Fields{"group": "updater websocket", "error": err}).Error("websocket client could not decrypt incoming packets")
				break
			}

			receiveChan <- messageDecrypted
		}
	}

	writeLoop := func() {
		defer func() {
			_ = client.Close()
		}()
		for {
			select {
			case message, ok := <-sendChan:
				if !ok {
					_ = client.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				err := client.WriteMessage(websocket.TextMessage, updater.sendCipher.Encrypt(nil, nil, message))
				if err != nil {
					updater.log.WithFields(logrus.Fields{"group": "updater websocket", "error": err}).Error("websocket could not write message")
				}

			case <-weHaveQuitChan:
				_ = client.WriteMessage(websocket.CloseMessage, []byte{})
				updater.log.Info("closing websocket connection")
				return
			}
		}
	}

	go readLoop()
	go writeLoop()

	return sendChan, weHaveQuitChan, receiveChan, remoteHasQuitChan
}

// initializeNoise sets up a new noise connection. First a fresh keypair is generated if none is locally found.
// Afterwards a XX handshake is performed. This is a three part handshake required to authenticate both parties.
// The resulting pairing code is then displayed to the user to check if it matches what is displayed on the other party's device.
func (updater *Updater) initializeNoise(client *websocket.Conn, bitboxBaseID string) error {
	cipherSuite := noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashSHA256)
	keypair := updater.configGetAppNoiseStaticKeypair()
	if keypair == nil {
		updater.log.Info("noise static keypair created")
		kp, err := cipherSuite.GenerateKeypair(rand.Reader)
		if err != nil {
			return errp.New("unable to generate a new noise keypair for the wallet app communication with the BitBox Base")
		}
		keypair = &kp
		if err := updater.configSetAppNoiseStaticKeypair(keypair); err != nil {
			updater.log.WithError(err).Error("could not store app noise static keypair")
			// Not a critical error, ignore.
		}
	}
	handshake, err := noise.NewHandshakeState(noise.Config{
		CipherSuite:   cipherSuite,
		Random:        rand.Reader,
		Pattern:       noise.HandshakeXX,
		StaticKeypair: *keypair,
		Prologue:      []byte("Noise_XX_25519_ChaChaPoly_SHA256"),
		Initiator:     true,
	})
	if err != nil {
		return errp.New("failed to generate a new noise handshake state for the wallet app communication with the BitBox Base")
	}

	//Ask the BitBox Base to begin the noise 'XX' handshake
	err = client.WriteMessage(1, []byte(opICanHasHandShaek))
	if err != nil {
		return errp.New("unable to write BitBox Base Handshake request to websocket")
	}
	_, responseBytes, err := client.ReadMessage()
	if err != nil {
		return errp.New("unable to read BitBox Base Handshake response from websocket")
	}
	if string(responseBytes) != string(responseSuccess) {
		return errp.New("no ACK received from BitBox Base as response to hanshake request")
	}

	// Do 3 part noise 'XX' handshake.
	msg, _, _, err := handshake.WriteMessage(nil, nil)
	if err != nil {
		return errp.New("noise failed to write the first handshake message")
	}
	err = client.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		return errp.New("the websocket failed writing the first noise handshake message")
	}
	_, responseBytes, err = client.ReadMessage()
	if err != nil {
		return errp.New("the websocket failed reading the second noise handshake message")
	}
	_, _, _, err = handshake.ReadMessage(nil, responseBytes)
	if err != nil {
		return errp.New("noise failed to read the second handshake message")
	}
	msg, updater.receiveCipher, updater.sendCipher, err = handshake.WriteMessage(nil, nil)
	if err != nil {
		return errp.New("noise failed to write the third handshake message")
	}
	err = client.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		return errp.New("the websocket failed writing the third handshake message")
	}

	// Check if the user already authenticated the channel binding hash
	_, responseBytes, err = client.ReadMessage()
	if err != nil {
		return errp.New("the websocket failed writing the pairingVerificationRequiredByBitBoxBase message at the verification stage of the noise handshake")
	}

	updater.bitboxBaseNoiseStaticPubkey = handshake.PeerStatic()
	if len(updater.bitboxBaseNoiseStaticPubkey) != 32 {
		return errp.New("expected 32 byte remote static pubkey")
	}

	pairingVerificationRequiredByApp := !updater.configContainsBitBoxBaseStaticPubkey(
		updater.bitboxBaseNoiseStaticPubkey)
	pairingVerificationRequiredByBase := string(responseBytes) == responseNeedsPairing

	// Do the user verification of the channel binding hash if either the app or base require it
	if pairingVerificationRequiredByBase || pairingVerificationRequiredByApp {
		channelHashBase32 := base32.StdEncoding.EncodeToString(handshake.ChannelBinding())
		updater.channelHash = fmt.Sprintf(
			"%s %s\n%s %s",
			channelHashBase32[:5],
			channelHashBase32[5:10],
			channelHashBase32[10:15],
			channelHashBase32[15:20])
		updater.Notify(observable.Event{
			Subject: fmt.Sprintf("/bitboxbases/%s/pairinghash", bitboxBaseID),
			Action:  action.Replace,
			Object:  updater.channelHash,
		})
		err = client.WriteMessage(websocket.BinaryMessage, []byte(opICanHasPairinVerificashun))
		if err != nil {
			return errp.New("the websocket failed writing the pairingVerificationRequiredByApp message at the verification stage of the noise handshake")
		}

		// Wait for the base to reply with responseSuccess, then proceed
		_, responseBytes, err := client.ReadMessage()
		if err != nil {
			return errp.New("websocket failed reading the pairing response from the BitBox Base")
		}
		updater.channelHashBitBoxBaseVerified = string(responseBytes) == string(responseSuccess)
		if updater.channelHashBitBoxBaseVerified {
			err = updater.configAddBitBoxBaseStaticPubkey(updater.bitboxBaseNoiseStaticPubkey)
			if err != nil {
				updater.log.Error("Pairing Successful, but unable to write bitboxBaseNoiseStaticPubkey to file")
			}
		} else {
			updater.sendCipher = nil
			updater.receiveCipher = nil
			updater.channelHash = ""
			return errp.New("pairing with BitBox Base failed")
		}
	}
	updater.channelHashAppVerified = true
	return nil
}
