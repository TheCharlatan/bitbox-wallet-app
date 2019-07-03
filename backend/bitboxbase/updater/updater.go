// Copyright 2019 Shift Devices AG
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

//Package updater manages the connection with the bitboxbase, establishing a websocket listener and sending events when receiving packets.
package updater

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	basemessages "github.com/digitalbitbox/bitbox-wallet-app/backend/bitboxbase/updater/messages"
	"github.com/digitalbitbox/bitbox-wallet-app/util/errp"
	"github.com/digitalbitbox/bitbox-wallet-app/util/logging"
	"github.com/digitalbitbox/bitbox-wallet-app/util/observable"
	"github.com/digitalbitbox/bitbox-wallet-app/util/observable/action"
	"github.com/flynn/noise"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"

	"github.com/sirupsen/logrus"
)

//go:generate protoc --go_out=import_path=basemessages:. messages/bbb.proto

const (
	opICanHasHandShaek          = "h"
	opICanHasPairinVerificashun = "v"
	responseSuccess             = "\x00"
	responseNeedsPairing        = "\x01"
)

// MiddlewareInfo holds some sample information from the BitBox Base
type MiddlewareInfo struct {
	Blocks         int64   `json:"blocks"`
	Difficulty     float64 `json:"difficulty"`
	LightningAlias string  `json:"lightningAlias"`
}

// Updater implements observable blockchainInfo.
type Updater struct {
	observable.Implementation
	middlewareInfo      *MiddlewareInfo
	log                 *logrus.Entry
	address             string
	bitboxBaseConfigDir string

	bitboxBaseNoiseStaticPubkey   []byte
	channelHash                   string
	channelHashAppVerified        bool
	channelHashBitBoxBaseVerified bool
	sendCipher, receiveCipher     *noise.CipherState

	weHaveQuitChan  chan<- struct{}
	bitboxBaseEvent chan proto.Message
	apiMap          map[string]chan *basemessages.BitBoxBaseOut

	onConnectionFailure func(string)
}

// MiddlewareInfo returns the last received blockchain information packet from the middleware
func (updater *Updater) MiddlewareInfo() *MiddlewareInfo {
	return updater.middlewareInfo
}

// NewUpdater returns a new bitboxbase updater.
func NewUpdater(address string, bitboxBaseConfigDir string, onConnectionFailure func(string)) *Updater {
	updater := &Updater{
		log:                 logging.Get().WithGroup("bitboxbase"),
		address:             address,
		middlewareInfo:      &MiddlewareInfo{},
		bitboxBaseConfigDir: bitboxBaseConfigDir,
		onConnectionFailure: onConnectionFailure,
	}
	return updater
}

// Connect starts the websocket go routine, first checking if the middleware is reachable,
// then establishing a websocket connection, then authenticating and encrypting all further traffic with noise.
func (updater *Updater) Connect(address string, bitboxBaseID string, bitboxBaseEvent chan proto.Message, apiMap map[string]chan *basemessages.BitBoxBaseOut) error {
	updater.bitboxBaseEvent = bitboxBaseEvent
	updater.apiMap = apiMap
	response, err := http.Get("http://" + address + "/")
	if err != nil {
		updater.log.Println("No response from middleware", err)
		return err
	}

	if response.StatusCode != http.StatusOK {
		updater.log.Println("Received http status code from middleware other than 200")
		return err
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		updater.log.Println("Body Bytes not read properly")
		return err
	}
	_, err = regexp.MatchString("OK!", string(bodyBytes))
	if err != nil {
		return errp.New("updater: Unexpected Response Body Bytes")
	}
	if err = response.Body.Close(); err != nil {
		return errp.New("updater: Failed to close Get Env response")
	}
	updater.log.Printf("connecting to base websocket")
	client, _, err := websocket.DefaultDialer.Dial("ws://"+updater.address+"/ws", nil)
	if err != nil {
		return errp.New("updater: failed to create new websocket client")
	}
	if err = updater.initializeNoise(client, bitboxBaseID); err != nil {
		return err
	}

	sendChan, weHaveQuitChan, receiveChan, remoteHasQuitChan := updater.runWebsocket(client)
	updater.weHaveQuitChan = weHaveQuitChan
	go listenWebsocket(sendChan, receiveChan, remoteHasQuitChan, updater, bitboxBaseID)
	return nil
}

//Stop shuts down the websocket connection with the base
func (updater *Updater) Stop() {
	close(updater.weHaveQuitChan)
}

func listenWebsocket(sendChan chan<- []byte, receiveChan <-chan []byte, quit <-chan struct{}, updater *Updater, bitboxBaseID string) {
	for {
		select {

		case message := <-receiveChan:
			incoming := &basemessages.BitBoxBaseOut{}
			if err := proto.Unmarshal(message, incoming); err != nil {
				updater.log.Println("protobuf unmarshal of incoming packet failed")
			}
			switch t := incoming.BitBoxBaseOut.(type) {
			case *basemessages.BitBoxBaseOut_BaseSystemEnvOut:
				updater.apiMap["systemEnv"] <- incoming
			case *basemessages.BitBoxBaseOut_BaseMiddlewareInfoOut:
				middlewareInfoIncoming := t //this is ugly, but golangci's gocritci wants it this way
				updater.middlewareInfo.LightningAlias = middlewareInfoIncoming.BaseMiddlewareInfoOut.LightningAlias
				updater.middlewareInfo.Blocks = middlewareInfoIncoming.BaseMiddlewareInfoOut.Blocks
				updater.middlewareInfo.Difficulty = float64(middlewareInfoIncoming.BaseMiddlewareInfoOut.Difficulty)
				updater.Notify(observable.Event{
					Subject: fmt.Sprintf("/bitboxbases/%s/middlewareinfo", bitboxBaseID),
					Action:  action.Replace,
					Object:  updater.middlewareInfo,
				})
				updater.log.Printf("Received blockinfo: %v , from id: %s", updater.middlewareInfo, bitboxBaseID)
			}

		case <-quit:
			updater.log.Error("Websocket closing, the BitBox Base " + bitboxBaseID + " disconnected")
			updater.onConnectionFailure(updater.address)
			return
		case event := <-updater.bitboxBaseEvent:
			data, err := proto.Marshal(event)
			if err != nil {
				updater.log.Println("Just print this generic error")
			}
			sendChan <- data
		}

	}
}
