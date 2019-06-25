package util

import (
	"fmt"
	log "github.com/op/go-logging"
	"github.com/spf13/viper"
	"net/http"
)

var logger = log.MustGetLogger("server/http")

var srv *http.Server

//deprecated
func StartHttpServer(vp *viper.Viper, h http.Handler) error {

	return StartHttpServerCustom(vp.GetString("host"), vp.GetInt("port"), h)
}

func StartHttpServerCustom(host string, port int, h http.Handler) error {
	if srv != nil {
		return fmt.Errorf("Server is running")
	}

	if port > 65535 {
		return fmt.Errorf("Invalid port: %d", port)
	} else if port == 0 {
		port = 8080
	}

	var listenaddr string
	if host == "" {
		listenaddr = fmt.Sprintf("%d", port)
	} else {
		listenaddr = fmt.Sprintf("%s:%d", host, port)
	}

	srv = &http.Server{Addr: listenaddr, Handler: h}

	// Start HTTP Server
	logger.Infof("Start HTTP Server at %s", listenaddr)
	err := srv.ListenAndServe()

	logger.Infof("Http Server is stopped: %v", err)
	return err

}

func IsHttpServerRunning() bool { return srv != nil }

func StopHttpServer() error {

	logger.Infof("Stop RPC server")

	if srv != nil {
		defer func() {
			srv = nil
		}()
		err := srv.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
