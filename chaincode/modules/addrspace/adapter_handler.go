package addrspace

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txh "hyperledger.abchain.org/chaincode/lib/txhandle"

	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
)

type regHandler struct{ AddrSpaceConfig }
type queryHandler struct{ AddrSpaceConfig }

func RegHandler(cfg AddrSpaceConfig) txh.TxHandler {
	return regHandler{cfg}
}

func QueryHandler(cfg AddrSpaceConfig) txh.TxHandler {
	return queryHandler{cfg}
}

func (h regHandler) Msg() proto.Message   { return new(empty.Empty) }
func (h queryHandler) Msg() proto.Message { return new(empty.Empty) }

func (h regHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	err := h.NewTx(stub, parser.GetNonce()).RegisterCC()
	if err != nil {
		return nil, err
	}

	return []byte("done"), nil
}

func (h queryHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	prefix, err := h.NewTx(stub, parser.GetNonce()).QueryPrefix()
	if err != nil {
		return nil, err
	}
	return rpc.EncodeRPCResult(&wrappers.BytesValue{Value: prefix})

}
