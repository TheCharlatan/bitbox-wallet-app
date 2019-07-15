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

package bitboxbase

import (
	"strings"
	"time"

	"github.com/digitalbitbox/bitbox-wallet-app/backend/bitboxbase/updater"
	basemessages "github.com/digitalbitbox/bitbox-wallet-app/backend/bitboxbase/updater/messages"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/btc/electrum"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/config"
	"github.com/digitalbitbox/bitbox-wallet-app/util/logging"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

// Interface represents bitbox base.
type Interface interface {
	// Identifier returns the bitboxBaseID.
	Identifier() string

	// GetUpdater returns the updater so we can listen to its events.
	GetUpdaterInstance() *updater.Updater

	// Close tells the bitboxbase to close all connections.
	Close()

	// GetRegisterTime implements a getter for the timestamp of when the bitboxBase was registered
	GetRegisterTime() time.Time

	// MiddlewareInfo returns some blockchain information.
	MiddlewareInfo() interface{}

	// ConnectElectrum connects to the electrs server on the base and configures the backend accordingly
	ConnectElectrum() error
}

// BitBoxBase provides the dictated bitboxbase api to communicate with the base
type BitBoxBase struct {
	bitboxBaseID             string //This is just the ip currently
	registerTime             time.Time
	address                  string
	updaterInstance          *updater.Updater
	electrsRPCPort           string
	network                  string
	log                      *logrus.Entry
	config                   *config.Config
	bitboxBaseConfigDir      string
	bitboxBaseEventsOutgoing chan proto.Message
	apiMap                   map[string]chan *basemessages.BitBoxBaseOut
}

//NewBitBoxBase creates a new bitboxBase instance
func NewBitBoxBase(address string, id string, config *config.Config, bitboxBaseConfigDir string, onConnectionFailure func(string)) (*BitBoxBase, error) {
	bitboxBase := &BitBoxBase{
		log:                      logging.Get().WithGroup("bitboxbase"),
		bitboxBaseID:             id,
		address:                  strings.Split(address, ":")[0],
		updaterInstance:          updater.NewUpdater(address, bitboxBaseConfigDir, onConnectionFailure),
		registerTime:             time.Now(),
		config:                   config,
		bitboxBaseConfigDir:      bitboxBaseConfigDir,
		bitboxBaseEventsOutgoing: make(chan proto.Message),
		apiMap:                   make(map[string]chan *basemessages.BitBoxBaseOut),
	}
	bitboxBase.apiMap["systemEnv"] = make(chan *basemessages.BitBoxBaseOut)
	err := bitboxBase.updaterInstance.Connect(address, bitboxBase.bitboxBaseID, bitboxBase.bitboxBaseEventsOutgoing, bitboxBase.apiMap)
	if err != nil {
		return nil, err
	}

	bitboxBase.getAndSetSystemEnv()
	return bitboxBase, err
}

func (base *BitBoxBase) getAndSetSystemEnv() {
	outgoing := &basemessages.BitBoxBaseIn{
		BitBoxBaseIn: &basemessages.BitBoxBaseIn_BaseSystemEnvIn{
			BaseSystemEnvIn: &basemessages.BaseSystemEnvIn{},
		},
	}
	base.bitboxBaseEventsOutgoing <- outgoing
	incoming := <-base.apiMap["systemEnv"]
	systemEnvIncoming, ok := incoming.BitBoxBaseOut.(*basemessages.BitBoxBaseOut_BaseSystemEnvOut)
	if ok {
		base.log.Println("Received SystemEnv information", systemEnvIncoming)
	}

	base.electrsRPCPort = systemEnvIncoming.BaseSystemEnvOut.GetElectrsRPCPort()
	base.network = systemEnvIncoming.BaseSystemEnvOut.GetElectrsRPCPort()
}

// ConnectElectrum connects to the electrs server on the base and configures the backend accordingly
func (base *BitBoxBase) ConnectElectrum() error {
	electrumAddress := base.address + ":" + base.electrsRPCPort

	electrumCert, err := electrum.DownloadCert(electrumAddress)
	if err != nil {
		base.log.WithField("ElectrumIP: ", electrumAddress).Error(err.Error())
		return err
	}

	if err := electrum.CheckElectrumServer(
		electrumAddress,
		electrumCert,
		base.log); err != nil {
		base.log.WithField("ElectrumIP: ", electrumAddress).Error(err.Error())
		return err
	}

	base.log.WithField("ElectrumAddress:", electrumAddress).Debug("Setting config to base electrum Server...")

	// BaseBtcConfig sets the TBTC configs to the provided cert and ip.
	if base.isTestnet() {
		base.config.SetTBTCElectrumServers(electrumAddress, electrumCert)
	} else {
		base.config.SetBTCElectrumServers(electrumAddress, electrumCert)
	}
	// Disable Litecoin and Ethereum accounts - we do not want any more traffic hitting other servers
	base.config.SetBtcOnly()

	if err := base.config.SetAppConfig(base.config.AppConfig()); err != nil {
		return err
	}
	return nil
}

// GetUpdaterInstance return ths current instance of the updater
func (base *BitBoxBase) GetUpdaterInstance() *updater.Updater {
	return base.updaterInstance
}

// MiddlewareInfo returns the received MiddlewareInfo packet from the updater
func (base *BitBoxBase) MiddlewareInfo() interface{} {
	return base.updaterInstance.MiddlewareInfo()
}

// Identifier implements a getter for the bitboxBase ID
func (base *BitBoxBase) Identifier() string {
	return base.bitboxBaseID
}

// GetRegisterTime implements a getter for the timestamp of when the bitboxBase was registered
func (base *BitBoxBase) GetRegisterTime() time.Time {
	return base.registerTime
}

// isTestnet returns a boolean that is true when connected to a base serving testnet and false otherwise
func (base *BitBoxBase) isTestnet() bool {
	return base.network == "testnet"
}

// Close implements a method to unset the bitboxBase
func (base *BitBoxBase) Close() {
	base.updaterInstance.Stop()
}
