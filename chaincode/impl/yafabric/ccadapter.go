package fabric_impl

import (
	"github.com/abchain/fabric/core/chaincode/shim"
	public_shim "hyperledger.abchain.org/chaincode/shim"
)

type adapter struct {
	public_shim.Chaincode
}

var nothingArg = [][]byte{}

func obtainByteArgs(stub shim.ChaincodeStubInterface) [][]byte {
	args := stub.GetArgs()
	if len(args) < 1 {
		return nothingArg
	} else {
		return args[1:]
	}
}

var logger = shim.NewLogger("aecc")

//we provide a DO-NOTHING deployment entry
func (a adapter) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Infof("A chaincode [%v] is deployed", a.Chaincode)
	return []byte("ok"), nil
}

func (a adapter) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return a.Chaincode.Invoke(stubAdapter{stub}, function, obtainByteArgs(stub), false)
}

func (a adapter) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return a.Chaincode.Invoke(stubAdapter{stub}, function, obtainByteArgs(stub), true)
}

func GenYAfabricCC(cc public_shim.Chaincode) shim.Chaincode {
	return adapter{cc}
}
