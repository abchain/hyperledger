package util

import (
	"encoding/json"

	"net/http"

	"github.com/gocraft/web"
	"hyperledger.abchain.org/core/utils"
)

//a null-base to provide more elastic
type FabricClientBase struct {
	debugData    interface{}
	RespWrapping func(interface{}) interface{}
}

func normalHeader(rw web.ResponseWriter) {

	// Set response content type
	rw.Header().Set("Content-Type", "application/json")

	// Enable CORS (default option handler will handle OPTION and set Access-Control-Allow-Method properly)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "accept, content-type")

	// Set response status ok
	rw.WriteHeader(http.StatusOK)

}

func (s *FabricClientBase) Normal(rw web.ResponseWriter, v interface{}) {

	normalHeader(rw)
	// Create response encoder
	if s.RespWrapping != nil {
		v = s.RespWrapping(v)
	}

	// logger.Debugf("Normal finish, output %v", v)

	json.NewEncoder(rw).Encode(utils.JRPCSuccess(v))
}

func (s *FabricClientBase) NormalError(rw web.ResponseWriter, e error) {

	normalHeader(rw)
	json.NewEncoder(rw).Encode(utils.JRPCError(e, s.debugData))
}

func (s *FabricClientBase) NormalErrorF(rw web.ResponseWriter, code int, message string) {

	normalHeader(rw)
	json.NewEncoder(rw).Encode(utils.JRPCErrorF(code, message, s.debugData))
}

func NormalEnd(rw web.ResponseWriter, v interface{}) {

	normalHeader(rw)
	json.NewEncoder(rw).Encode(utils.JRPCSuccess(v))
}

func NormalError(rw web.ResponseWriter, e error) {

	normalHeader(rw)
	json.NewEncoder(rw).Encode(utils.JRPCError(e, nil))
}

func NormalErrorF(rw web.ResponseWriter, code int, message string) {

	normalHeader(rw)
	json.NewEncoder(rw).Encode(utils.JRPCErrorF(code, message, nil))
}
