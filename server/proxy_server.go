package server

import (
	"ce102/config"
	"ce102/controller"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	controller controller.Controller
}

func (s *Server) addIoTDevice(w http.ResponseWriter, r *http.Request) {
	log.Println("handler addIoTDevice")
	defer r.Body.Close()

	var iotDev config.IotConfig
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&iotDev); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := s.controller.AddIoTDevice(iotDev); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) rmIoTDevice(w http.ResponseWriter, r *http.Request) {
	log.Println("handler rmIoTDevice")
	defer r.Body.Close()

	deviceNames := r.URL.Query()["deviceName"]
	if len(deviceNames) == 0 {
		log.Println("device name not found")
		fmt.Fprintf(w, "set device name")
		return
	}
	deviceName := deviceNames[0]

	if err := s.controller.RmIoTDevice(deviceName); err != nil {
		log.Println(err)
		fmt.Fprintf(w, "wrong device name")
		return
	}
}

func (s *Server) stopObserveDevice(w http.ResponseWriter, r *http.Request) {
	log.Println("handler stopObserveDevice")
	defer r.Body.Close()

	deviceNames := r.URL.Query()["deviceName"]
	if len(deviceNames) == 0 {
		log.Println("device name not found")
		fmt.Fprintf(w, "set device name")
		return
	}
	deviceName := deviceNames[0]

	if err := s.controller.StopObserveDevice(deviceName); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) getInformationFromIotDevice(w http.ResponseWriter, r *http.Request) {
	log.Println("handler getInformationFromIotDevice")
	defer r.Body.Close()
	deviceNames := r.URL.Query()["deviceName"]
	if len(deviceNames) == 0 {
		log.Println("device name not found")
		fmt.Fprintf(w, "set device name")
		return
	}
	deviceName := deviceNames[0]

	inf, err := s.controller.GetInformation(deviceName)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(inf)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) getLogs(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Println("handler getLogs")
	countLogsArr := r.URL.Query()["countLogs"]
	if len(countLogsArr) == 0 {
		log.Println("count logs not found")
		fmt.Fprintf(w, "set count logs")
		return
	}
	countLogsStr := countLogsArr[0]
	countLogs, err := strconv.Atoi(countLogsStr)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logs, err := s.controller.GetLastNRowsLogs(countLogs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allLogs := strings.Join(logs, "\n")
	_, err = w.Write([]byte(allLogs))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) StartServer(config config.Config, controller controller.Controller) {
	s.controller = controller

	http.HandleFunc("/device/metrics", s.getInformationFromIotDevice)
	http.HandleFunc("/logs", s.getLogs)
	http.HandleFunc("/device/add", s.addIoTDevice)
	http.HandleFunc("/device/rm", s.rmIoTDevice)
	http.HandleFunc("/device/observer/stop", s.stopObserveDevice)

	fmt.Println("Server is listening... ", config.ProxyServerAddr)
	log.Fatal(http.ListenAndServe(config.ProxyServerAddr, nil))
}
