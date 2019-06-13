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
	"encoding/json"
	"strings"
	"time"

	"github.com/digitalbitbox/bitbox-wallet-app/backend/bitboxbase/updater"
	"github.com/digitalbitbox/bitbox-wallet-app/util/logging"
	"github.com/sirupsen/logrus"
)

// Interface represents bitbox base.
type Interface interface {
	Init(testing bool)

	// Identifier returns the bitboxBaseID.
	Identifier() string

	// GetUpdater returns the updater so we can listen to its events.
	GetUpdaterInstance() *updater.Updater

	// Close tells the bitboxbase to close all connections.
	Close()

	//GetRegisterTime implements a getter for the timestamp of when the bitboxBase was registered
	GetRegisterTime() time.Time

	// BlockInfo returns some blockchain information.
	BlockInfo() interface{}

	// GetIP implement a getter for the IP under which the base is reachable
	GetIP() string

	// GetElectrsRPCPort implements a getter for the electrs rpc port
	GetElectrsRPCPort() string

	// GetNetwork implements a getter for the network type, either mainnet or testnet
	GetNetwork() string

    // GetElectrsConnected implements a getter for a boolean indicating wether this base is connected to electrs
    GetElectrsConected() bool
}

//BitBoxBase provides the dictated bitboxbase api to communicate with the base
type BitBoxBase struct {
	bitboxBaseID    string //This is just the ip at the moment, but will be an actual unique string, once the noise pairing is implemented
	registerTime    time.Time
	ip              string
	closed          bool
	updaterInstance *updater.Updater
	electrsRPCPort  string
    ElectrsConnected bool
	network         string
	log             *logrus.Entry
}

//NewBitBoxBase creates a new bitboxBase instance
func NewBitBoxBase(ip string, id string) (*BitBoxBase, error) {
	bitboxBase := &BitBoxBase{
		log:             logging.Get().WithGroup("bitboxbase"),
		bitboxBaseID:    id,
		closed:          false,
        ElectrsConnected: false,
		ip:              strings.Split(ip, ":")[0],
		updaterInstance: updater.NewUpdater(ip),
		registerTime:    time.Now(),
	}
	err := bitboxBase.GetUpdaterInstance().Connect(ip, bitboxBase.bitboxBaseID)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := bitboxBase.GetUpdaterInstance().GetEnv()
	if err != nil {
		return nil, err
	}
	var envData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &envData); err != nil {
		bitboxBase.log.WithError(err).Error(" Failed to unmarshal GetEnv body bytes")
		bitboxBase.GetUpdaterInstance().Stop()
		return nil, err
	}
	bitboxBase.electrsRPCPort, _ = envData["electrsRPCPort"].(string)
	bitboxBase.network, _ = envData["network"].(string)

	return bitboxBase, err
}

//GetUpdaterInstance return ths current instance of the updater
func (base *BitBoxBase) GetUpdaterInstance() *updater.Updater {
	return base.updaterInstance
}

//BlockInfo returns the received blockinfo packet from the updater
func (base *BitBoxBase) BlockInfo() interface{} {
	return base.GetUpdaterInstance().MiddlewareInfo()
}

//GetIP implements a getter for the bitboxBase ip
func (base *BitBoxBase) GetIP() string {
	return base.ip
}

//Identifier implements a getter for the bitboxBase ID
func (base *BitBoxBase) Identifier() string {
	return base.bitboxBaseID
}

// GetElectrsRPCPort implements a getter for the electrsport
func (base *BitBoxBase) GetElectrsRPCPort() string {
	return base.electrsRPCPort
}

//GetNetwork implements a getter for the network string (either mainnet or testnet)
func (base *BitBoxBase) GetNetwork() string {
	return base.network
}

//GetRegisterTime implements a getter for the timestamp of when the bitboxBase was registered
func (base *BitBoxBase) GetRegisterTime() time.Time {
	return base.registerTime
}

//GetElectrsConnected is a getter for the boolean indicating if the base is connected with electrs
func (base *BitBoxBase) GetElectrsConencted() bool {
    return base.ElectrsConnected
}

//Close implements a method to unset the bitboxBase
func (base *BitBoxBase) Close() {
	base.GetUpdaterInstance().Stop()
	base.closed = true
}

//Init initializes the bitboxBase
func (base *BitBoxBase) Init(testing bool) {
}
