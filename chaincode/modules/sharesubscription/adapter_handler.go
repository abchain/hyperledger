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

	contract := make(map[string]int32)
	for _, m := range msg.ContractBody {
		addr, err := txutil.NewAddressFromPBMessage(m.Addr)
		if err != nil {
			return nil, err
		}

		contract[addr.ToString()] = m.Weight
	}
	return h.NewTx(stub, parser.GetNounce()).New(contract, msg.DelegatorAddr.GetHash())
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

	resp, err := h.NewTx(stub, parser.GetNounce()).Redeem(contract.Hash, toAmount(msg.Amount), redeemTo)
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

	err, data := h.NewTx(stub, parser.GetNounce()).Query(contAddr.Hash)
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

	err, data := h.NewTx(stub, parser.GetNounce()).QueryOne(contAddr.Hash, member.Hash)
	if err != nil {
		return nil, err
	}

	return rpc.EncodeRPCResult(data.ToPB())
}
