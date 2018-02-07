package service

import (
	"fmt"
	log "github.com/abchain/fabric/peerex/logging"
	"net/http"
)

var logger = log.InitLogger("server")

var srv *http.Server

func startHttpServer(host string, port int) error {

	listenaddr := fmt.Sprintf("%s:%d", host, port)

	router := buildRouter()
	srv = &http.Server{Addr: listenaddr, Handler: router}

	// Start HTTP Server
	logger.Infof("Start RPC Server at http://%s", listenaddr)
	err := srv.ListenAndServe()

	logger.Infof("The RPC Server is stopped: %v", err)
	return err
}

func stopHttpServer() error {

	logger.Infof("Stop RPC server")

	if srv != nil {
		err := srv.Shutdown(nil)
		if err != nil {
			return err
		}
	}

	return nil
}
