package subscription

import (
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/lib/caller"
	pb "hyperledger.abchain.org/chaincode/modules/sharesubscription/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
)

type newContractHandler struct{ ContractConfig }
type redeemHandler struct{ ContractConfig }
type queryHandler struct{ ContractConfig }
type memberQueryHandler struct{ ContractConfig }

func NewContractHandler(cfg ContractConfig) newContractHandler {
	return newContractHandler{ContractConfig: cfg}
}

func RedeemHandler(cfg ContractConfig) redeemHandler {
	return redeemHandler{ContractConfig: cfg}
}

func QueryHandler(cfg ContractConfig) queryHandler {
	return queryHandler{ContractConfig: cfg}
}
func MemberQueryHandler(cfg ContractConfig) memberQueryHandler {
	return memberQueryHandler{ContractConfig: cfg}
}

func (h newContractHandler) Msg() proto.Message { return new(pb.RegContract) }
func (h redeemHandler) Msg() proto.Message      { return new(pb.RedeemContract) }
func (h queryHandler) Msg() proto.Message       { return new(pb.QueryContract) }
func (h memberQueryHandler) Msg() proto.Message { return new(pb.QueryContract) }

func (h newContractHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	msg := parser.GetMessage().(*pb.RegContract)

	var addrs [][]byte
	var ratios []int

	for _, m := range msg.ContractBody {
		addr, err := txutil.NewAddressFromPBMessage(m.Addr)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr.Internal())
		ratios = append(ratios, int(m.GetWeight()))
	}
	return h.NewTx(stub, parser.GetNonce()).New_C(addrs, ratios)
}

func (h redeemHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := parser.GetMessage().(*pb.RedeemContract)

	contract, err := txutil.NewAddressFromPBMessage(msg.Contract)
	if err != nil {
		return nil, err
	}

	var redeemTo [][]byte

	for _, addr := range msg.Redeems {
		redeemTo = append(redeemTo, addr.GetHash())
	}

	resp, err := h.NewTx(stub, parser.GetNonce()).Redeem_C(contract.Hash, toAmount(msg.Amount), redeemTo)
	if err != nil {
		return nil, err
	}

	return rpc.EncodeRPCResult(resp)
}

func (h queryHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	msg := parser.GetMessage().(*pb.QueryContract)

	contAddr, err := txutil.NewAddressFromPBMessage(msg.ContractAddr)
	if err != nil {
		return nil, err
	}

	err, data := h.NewTx(stub, parser.GetNonce()).Query_C(contAddr.Hash)
	if err != nil {
		return nil, err
	}

	return rpc.EncodeRPCResult(data.ToPB())

}

func (h memberQueryHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := parser.GetMessage().(*pb.QueryContract)

	contAddr, err := txutil.NewAddressFromPBMessage(msg.ContractAddr)
	if err != nil {
		return nil, err
	}

	member, err := txutil.NewAddressFromPBMessage(msg.MemberAddr)
	if err != nil {
		return nil, err
	}

	err, data := h.NewTx(stub, parser.GetNonce()).QueryOne_C(contAddr.Hash, member.Hash)
	if err != nil {
		return nil, err
	}

	return rpc.EncodeRPCResult(data.ToPB())
}
