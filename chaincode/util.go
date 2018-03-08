package ccutil

import (
	"github.com/abchain/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("ccutil")

func ExecuteCC(cc shim.Chaincode) chan error {

	ret := make(chan error)

	go func() {
		err := shim.Start(cc)
		if err != nil {
			logger.Errorf("Error starting chaincode: %s", err)
		}

		ret <- err
	}()

	return ret
}
