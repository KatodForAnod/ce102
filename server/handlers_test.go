package server

import (
	"ce102/config"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	proxyServer Server
)

const serverAddr = "127.0.0.1:8080"

type Controller struct {
	ioTs []config.IotConfig
}

func (c *Controller) AddIoTDevice(device config.IotConfig) error {
	panic("implement me")
}

func (c *Controller) RmIoTDevice(deviceName string) error {
	panic("implement me")
}

func (c *Controller) StopObserveDevice(deviceName string) error {
	panic("implement me")
}

func (c *Controller) GetLastNRowsLogs(nRows int) ([]string, error) {
	if nRows < 0 {
		return []string{}, errors.New("wrong count rows")
	}
	return []string{"1 row", "2 row"}, nil
}

func (c *Controller) RemoveIoTDeviceObserve(ioTsConfig []config.IotConfig) error {
	return nil
}

func (c *Controller) NewIotDeviceObserve(iotConfig config.IotConfig) error {
	c.ioTs = append(c.ioTs, iotConfig)
	return nil
}

func (c *Controller) GetInformation(deviceName string) ([]byte, error) {
	for i, t := range c.ioTs {
		if t.DeviceName == deviceName {
			return []byte{byte(i)}, nil
		}
	}

	return []byte{}, errors.New("not found")
}

func Init() {
	controller := Controller{}
	proxyServer.StartServer(config.Config{ProxyServerAddr: serverAddr}, &controller)
}

func TestServer_addIotDevice(t *testing.T) {
	go Init()
	req := httptest.NewRequest(http.MethodGet, "/device/add?deviceName=testName&deviceAddr=:5600", nil)
	w := httptest.NewRecorder()
	proxyServer.getLogs(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestServer_getInformationFromIotDevice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/device/metrics?deviceName=testName", nil)
	w := httptest.NewRecorder()
	proxyServer.getLogs(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestServer_removeIotDevice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/device/rm?deviceName=testName", nil)
	w := httptest.NewRecorder()
	proxyServer.getLogs(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestServer_getLogs(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/logs?countLogs=2", nil)
	w := httptest.NewRecorder()
	proxyServer.getLogs(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	out := w.Body.String()
	outArr := strings.Split(out, "\n")
	if len(outArr) < 2 {
		t.FailNow()
	}
}
