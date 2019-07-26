package multisign

import (
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/lib/caller"
	pb "hyperledger.abchain.org/chaincode/modules/multisign/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
	txpb "hyperledger.abchain.org/protos"
)

type contractHandler struct{ MultiSignConfig }
type updateHandler struct{ MultiSignConfig }
type queryHandler struct{ MultiSignConfig }

func ContractHandler(cfg MultiSignConfig) contractHandler {
	return contractHandler{MultiSignConfig: cfg}
}
func UpdateHandler(cfg MultiSignConfig) updateHandler {
	return updateHandler{MultiSignConfig: cfg}
}
func QueryHandler(cfg MultiSignConfig) queryHandler {
	return queryHandler{MultiSignConfig: cfg}
}

func (h contractHandler) Msg() proto.Message { return new(pb.Contract) }
func (h updateHandler) Msg() proto.Message   { return new(pb.Update) }
func (h queryHandler) Msg() proto.Message    { return new(txpb.TxAddr) }

func (h contractHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	msg := parser.GetMessage().(*pb.Contract)
	var addrs [][]byte
	var weights []int32

	for _, m := range msg.Addrs {
		addrs = append(addrs, m.GetAddr().GetHash())
		weights = append(weights, m.GetWeight())
	}
	return h.NewTx(stub, parser.GetNonce()).Contract_C(msg.Threshold, addrs, weights)
}

func (h updateHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := parser.GetMessage().(*pb.Update)

	err := h.NewTx(stub, parser.GetNonce()).Update_C(msg.GetAddr().GetHash(),
		msg.GetFrom().GetHash(), msg.GetTo().GetHash())
	if err != nil {
		return nil, err
	}
	return []byte("done"), nil
}

func (h queryHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	msg := parser.GetMessage().(*txpb.TxAddr)

	err, data := h.NewTx(stub, parser.GetNonce()).Query_C(msg.GetHash())
	if err != nil {
		return nil, err
	}

	return rpc.EncodeRPCResult(data.ToPB())

}
