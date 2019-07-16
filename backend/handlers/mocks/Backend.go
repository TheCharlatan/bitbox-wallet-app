// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import accounts "github.com/digitalbitbox/bitbox-wallet-app/backend/accounts"
import bitboxbase "github.com/digitalbitbox/bitbox-wallet-app/backend/bitboxbase"
import coin "github.com/digitalbitbox/bitbox-wallet-app/backend/coins/coin"
import config "github.com/digitalbitbox/bitbox-wallet-app/backend/config"
import device "github.com/digitalbitbox/bitbox-wallet-app/backend/devices/device"

import keystore "github.com/digitalbitbox/bitbox-wallet-app/backend/keystore"
import language "golang.org/x/text/language"
import mock "github.com/stretchr/testify/mock"
import signing "github.com/digitalbitbox/bitbox-wallet-app/backend/signing"

// Backend is an autogenerated mock type for the Backend type
type Backend struct {
	mock.Mock
}

// Accounts provides a mock function with given fields:
func (_m *Backend) Accounts() []accounts.Interface {
	ret := _m.Called()

	var r0 []accounts.Interface
	if rf, ok := ret.Get(0).(func() []accounts.Interface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]accounts.Interface)
		}
	}

	return r0
}

// AccountsStatus provides a mock function with given fields:
func (_m *Backend) AccountsStatus() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// BitBoxBaseDeregister provides a mock function with given fields: bitboxBaseID
func (_m *Backend) BitBoxBaseDeregister(bitboxBaseID string) {
	_m.Called(bitboxBaseID)
}

// BitBoxBasesRegistered provides a mock function with given fields:
func (_m *Backend) BitBoxBasesRegistered() map[string]bitboxbase.Interface {
	ret := _m.Called()

	var r0 map[string]bitboxbase.Interface
	if rf, ok := ret.Get(0).(func() map[string]bitboxbase.Interface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]bitboxbase.Interface)
		}
	}

	return r0
}

// CheckElectrumServer provides a mock function with given fields: _a0, _a1
func (_m *Backend) CheckElectrumServer(_a0 string, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Coin provides a mock function with given fields: _a0
func (_m *Backend) Coin(_a0 string) (coin.Coin, error) {
	ret := _m.Called(_a0)

	var r0 coin.Coin
	if rf, ok := ret.Get(0).(func(string) coin.Coin); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coin.Coin)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Config provides a mock function with given fields:
func (_m *Backend) Config() *config.Config {
	ret := _m.Called()

	var r0 *config.Config
	if rf, ok := ret.Get(0).(func() *config.Config); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.Config)
		}
	}

	return r0
}

// CreateAndAddAccount provides a mock function with given fields: _a0, code, name, getSigningConfiguration, persist
func (_m *Backend) CreateAndAddAccount(_a0 coin.Coin, code string, name string, getSigningConfiguration func() (*signing.Configuration, error), persist bool) error {
	ret := _m.Called(_a0, code, name, getSigningConfiguration, persist)

	var r0 error
	if rf, ok := ret.Get(0).(func(coin.Coin, string, string, func() (*signing.Configuration, error), bool) error); ok {
		r0 = rf(_a0, code, name, getSigningConfiguration, persist)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DefaultAppConfig provides a mock function with given fields:
func (_m *Backend) DefaultAppConfig() config.AppConfig {
	ret := _m.Called()

	var r0 config.AppConfig
	if rf, ok := ret.Get(0).(func() config.AppConfig); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(config.AppConfig)
	}

	return r0
}

// Deregister provides a mock function with given fields: deviceID
func (_m *Backend) Deregister(deviceID string) {
	_m.Called(deviceID)
}

// DeregisterKeystore provides a mock function with given fields:
func (_m *Backend) DeregisterKeystore() {
	_m.Called()
}

// DevicesRegistered provides a mock function with given fields:
func (_m *Backend) DevicesRegistered() map[string]device.Interface {
	ret := _m.Called()

	var r0 map[string]device.Interface
	if rf, ok := ret.Get(0).(func() map[string]device.Interface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]device.Interface)
		}
	}

	return r0
}

// DownloadCert provides a mock function with given fields: _a0
func (_m *Backend) DownloadCert(_a0 string) (string, error) {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NotifyUser provides a mock function with given fields: _a0
func (_m *Backend) NotifyUser(_a0 string) {
	_m.Called(_a0)
}

// OnAccountInit provides a mock function with given fields: f
func (_m *Backend) OnAccountInit(f func(accounts.Interface)) {
	_m.Called(f)
}

// OnAccountUninit provides a mock function with given fields: f
func (_m *Backend) OnAccountUninit(f func(accounts.Interface)) {
	_m.Called(f)
}

// OnBitBoxBaseInit provides a mock function with given fields: f
func (_m *Backend) OnBitBoxBaseInit(f func(bitboxbase.Interface)) {
	_m.Called(f)
}

// OnBitBoxBaseUninit provides a mock function with given fields: f
func (_m *Backend) OnBitBoxBaseUninit(f func(string)) {
	_m.Called(f)
}

// OnDeviceInit provides a mock function with given fields: f
func (_m *Backend) OnDeviceInit(f func(device.Interface)) {
	_m.Called(f)
}

// OnDeviceUninit provides a mock function with given fields: f
func (_m *Backend) OnDeviceUninit(f func(string)) {
	_m.Called(f)
}

// Rates provides a mock function with given fields:
func (_m *Backend) Rates() map[string]map[string]float64 {
	ret := _m.Called()

	var r0 map[string]map[string]float64
	if rf, ok := ret.Get(0).(func() map[string]map[string]float64); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]map[string]float64)
		}
	}

	return r0
}

// Register provides a mock function with given fields: _a0
func (_m *Backend) Register(_a0 device.Interface) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(device.Interface) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegisterKeystore provides a mock function with given fields: _a0
func (_m *Backend) RegisterKeystore(_a0 keystore.Keystore) {
	_m.Called(_a0)
}

// RegisterTestKeystore provides a mock function with given fields: _a0
func (_m *Backend) RegisterTestKeystore(_a0 string) {
	_m.Called(_a0)
}

// Start provides a mock function with given fields:
func (_m *Backend) Start() <-chan interface{} {
	ret := _m.Called()

	var r0 <-chan interface{}
	if rf, ok := ret.Get(0).(func() <-chan interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan interface{})
		}
	}

	return r0
}

// SystemOpen provides a mock function with given fields: _a0
func (_m *Backend) SystemOpen(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Testing provides a mock function with given fields:
func (_m *Backend) Testing() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// TryMakeNewBase provides a mock function with given fields: ip
func (_m *Backend) TryMakeNewBase(ip string) (bool, error) {
	ret := _m.Called(ip)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(ip)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(ip)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserLanguage provides a mock function with given fields:
func (_m *Backend) UserLanguage() language.Tag {
	ret := _m.Called()

	var r0 language.Tag
	if rf, ok := ret.Get(0).(func() language.Tag); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(language.Tag)
	}

	return r0
}
