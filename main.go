package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/SwingbyProtocol/node-agent/config"
	"github.com/SwingbyProtocol/node-agent/types"
	"github.com/SwingbyProtocol/tx-indexer/api"
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			//_, filename := path.Split(f.File)
			paddedFuncname := fmt.Sprintf(" %-20v", funcname+"()")
			//paddedFilename := fmt.Sprintf("%17v", filename)
			return paddedFuncname, ""
		},
	})
	log.SetOutput(os.Stdout)
}

func main() {
	conf, err := config.NewDefaultConfig()
	if err != nil {
		log.Fatal(err)
	}
	apiConfig := &api.APIConfig{
		ListenREST: conf.RESTConfig.ListenAddr,
		ListenWS:   conf.WSConfig.ListenAddr,
		Actions:    []*api.Action{},
	}

	getStatus := func(w rest.ResponseWriter, r *rest.Request) {

		file, _ := ioutil.ReadFile("./data/node_status.json")

		data := types.NodeStatus{}

		_ = json.Unmarshal([]byte(file), &data)

		w.WriteJson(data)
	}
	getStatusAction := api.NewGet("/api/v1/status", getStatus)
	apiConfig.Actions = append(apiConfig.Actions, getStatusAction)

	// Create api server
	apiServer := api.NewAPI(apiConfig)
	// Start server
	apiServer.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGSTOP)
	signal := <-c
	// Backup operation
	log.Info(signal)
}
