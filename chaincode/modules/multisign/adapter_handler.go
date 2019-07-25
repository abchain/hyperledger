package multisign

import (
	"errors"
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

	addr2Weight := make(map[string]int32)
	for _, m := range msg.Addrs {
		addr, err := txutil.NewAddressFromPBMessage(m.GetAddr())
		if err != nil {
			return nil, err
		}

		addrs := addr.ToString()

		if _, existed := addr2Weight[addrs]; existed {
			return nil, errors.New("Duplicated contract address")
		}

		addr2Weight[addrs] = m.GetWeight()
	}
	return h.NewTx(stub, parser.GetNonce()).Contract(msg.Threshold, addr2Weight)
}

func (h updateHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := parser.GetMessage().(*pb.Update)

	var addrs [3]string
	var target []*txpb.TxAddr
	if msg.GetTo() == nil {
		target = []*txpb.TxAddr{msg.GetAddr(), msg.GetFrom()}
	} else {
		target = []*txpb.TxAddr{msg.GetAddr(), msg.GetFrom(), msg.GetTo()}
	}

	for i, pbAddr := range target {

		addr, err := txutil.NewAddressFromPBMessage(pbAddr)
		if err != nil {
			return nil, err
		}
		addrs[i] = addr.ToString()
	}

	err := h.NewTx(stub, parser.GetNonce()).Update(addrs[0], addrs[1], addrs[2])
	if err != nil {
		return nil, err
	}
	return []byte("done"), nil
}

func (h queryHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	msg := parser.GetMessage().(*txpb.TxAddr)

	accAddr, err := txutil.NewAddressFromPBMessage(msg)
	if err != nil {
		return nil, err
	}

	err, data := h.NewTx(stub, parser.GetNonce()).Query(accAddr.ToString())
	if err != nil {
		return nil, err
	}

	return rpc.EncodeRPCResult(data.ToPB())

}
