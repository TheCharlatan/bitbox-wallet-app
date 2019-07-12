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

package handlers_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/digitalbitbox/bitbox-wallet-app/backend"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/arguments"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/devices/usb"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/handlers"
	"github.com/digitalbitbox/bitbox-wallet-app/util/logging"
	"github.com/digitalbitbox/bitbox-wallet-app/util/system"
	"github.com/digitalbitbox/bitbox-wallet-app/util/test"
	"github.com/gorilla/mux"
)

// webdevEnvironment implements backend.Environment
type webdevEnvironment struct {
}

// DeviceInfos implements backend.Environment
func (webdevEnvironment) DeviceInfos() []usb.DeviceInfo {
	return usb.DeviceInfos()
}

// SystemOpen implements backend.Environment
func (webdevEnvironment) SystemOpen(url string) error {
	return system.Open(url)
}

// NotifyUser implements backend.Environment
func (webdevEnvironment) NotifyUser(text string) {
	log := logging.Get().WithGroup("servewallet")
	log.Infof("NotifyUser: %s", text)
	// We use system notifications on unix/macOS, the primary dev environments.
	switch runtime.GOOS {
	case "darwin":
		// #nosec G204
		err := exec.Command("osascript", "-e",
			fmt.Sprintf(`display notification "%s" with title \"BitBox Wallet DEV\"`, text))
		if err != nil {
			log.Error(err)
		}
	case "linux":
		// #nosec G204b
		err := exec.Command("notify-send", "BitBox Wallet DEV", text).Run()
		if err != nil {
			log.Error(err)
		}
	}
}

// List all routes with `go test backend/handlers/handlers_test.go -v`.
func TestListRoutes(t *testing.T) {
	connectionData := handlers.NewConnectionData(-1, "")
	backend, err := backend.NewBackend(arguments.NewArguments(
		test.TstTempDir("bitbox-wallet-listroutes-"), false, false, false, false, false),
		webdevEnvironment{},
	)
	if err != nil {
		fmt.Println(err)
	}
	handlers := handlers.NewHandlers(backend, connectionData)
	err = handlers.Router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			return err
		}
		if len(methods) == 0 {
			fmt.Println()
		}
		fmt.Print(pathTemplate)
		if len(methods) > 0 {
			fmt.Print(" (" + strings.Join(methods, ",") + ")")
		}
		/* The following methods are only available in a newer version of mux: */
		// queriesTemplates, err := route.GetQueriesTemplates()
		// if err == nil {
		// 	   fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		// }
		// queriesRegexps, err := route.GetQueriesRegexp()
		// if err == nil {
		// 	   fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		// }
		fmt.Println()
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func TestTestingHandler(t *testing.T) {
	connectionData := handlers.NewConnectionData(-1, "")

	backend, err := backend.NewBackend(arguments.NewArguments(
		test.TstTempDir("bitbox-wallet-listroutes-"), false, false, false, false, false),
		webdevEnvironment{},
	)
	if err != nil {
		fmt.Println(err)
	}
	handlers := handlers.NewHandlers(backend, connectionData)
	req, err := http.NewRequest("GET", "/api/testing", nil)
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
	rr := httptest.NewRecorder()
	handlers.Router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d. Got %d", http.StatusOK, status)
	}
}

func TestVersionHandler(t *testing.T) {
	connectionData := handlers.NewConnectionData(-1, "")

	backend, err := backend.NewBackend(arguments.NewArguments(
		test.TstTempDir("bitbox-wallet-listroutes-"), false, false, false, false, false),
		webdevEnvironment{},
	)
	if err != nil {
		fmt.Println(err)
	}
	handlers := handlers.NewHandlers(backend, connectionData)
	req, err := http.NewRequest("GET", "/api/version", nil)
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
	rr := httptest.NewRecorder()
	handlers.Router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d. Got %d", http.StatusOK, status)
	}
}

func TestUpdateHandler(t *testing.T) {
	connectionData := handlers.NewConnectionData(-1, "")

	backend, err := backend.NewBackend(arguments.NewArguments(
		test.TstTempDir("bitbox-wallet-listroutes-"), false, false, false, false, false),
		webdevEnvironment{},
	)
	if err != nil {
		fmt.Println(err)
	}
	handlers := handlers.NewHandlers(backend, connectionData)
	req, err := http.NewRequest("GET", "/api/update", nil)
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
	rr := httptest.NewRecorder()
	handlers.Router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d. Got %d", http.StatusOK, status)
	}
}

func TestConfigDefaultHandler(t *testing.T) {
	connectionData := handlers.NewConnectionData(-1, "")

	backend, err := backend.NewBackend(arguments.NewArguments(
		test.TstTempDir("bitbox-wallet-listroutes-"), false, false, false, false, false),
		webdevEnvironment{},
	)
	if err != nil {
		fmt.Println(err)
	}
	handlers := handlers.NewHandlers(backend, connectionData)
	req, err := http.NewRequest("GET", "/api/config/default", nil)
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
	rr := httptest.NewRecorder()
	handlers.Router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d. Got %d", http.StatusOK, status)
	}

	req, err = http.NewRequest("POST", "/api/config", bytes.NewBuffer([]byte(rr.Body.String())))
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
	rr = httptest.NewRecorder()
	handlers.Router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d. Got %d", http.StatusOK, status)
	}
}
