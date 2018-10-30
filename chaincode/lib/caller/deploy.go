package rpc

import (
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/shim"
)

type DeployHandler interface {
	Call(stub shim.ChaincodeStubInterface, deployarg []byte) error
}

func DeployCC(stub shim.ChaincodeStubInterface, args []string,
	handlers map[string]DeployHandler) error {

	//build args mapper
	var index string
	argmap := make(map[string]string)

	for _, v := range args {
		if index == "" {
			index = v
		} else {
			argmap[index] = v
			index = ""
		}
	}

	for k, h := range handlers {
		arg, ok := argmap[k]
		if !ok {
			return fmt.Errorf("No correspoinding arguments for deploying %s", k)
		}

		bt, err := base64.RawStdEncoding.DecodeString(arg)
		if err != nil {
			return fmt.Errorf("Deploy fail while decoding arg for %s: %s", k, err.Error())
		}

		err = h.Call(stub, bt)
		if err != nil {
			return fmt.Errorf("Deploy fail for handling %s: %s", k, err.Error())
		}
	}

	return nil
}

func BuildDeployArg(method string, msg proto.Message, args []string) ([]string, error) {
	bt, err := EncodeRPCResult(msg)
	if err != nil {
		return args, err
	}

	return append(args, method, base64.RawStdEncoding.EncodeToString(bt)), nil
}
