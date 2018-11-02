package subscription

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/lib/caller"
	pb "hyperledger.abchain.org/chaincode/modules/sharesubscription/protos"
	"hyperledger.abchain.org/chaincode/shim"
	"hyperledger.abchain.org/core/crypto"
	txutil "hyperledger.abchain.org/core/tx"
)

type RedeemMsg struct {
	msg        pb.RedeemContract
	redeemAddr *txutil.Address
}

type RegContractMsg struct {
	msg pb.RegContract
}

type newContractHandler struct {
	RegContractMsg
	ContractConfig
	pk *crypto.PublicKey
}

type redeemHandler struct {
	RedeemMsg
	ContractConfig
}

type queryHandler struct {
	msg pb.QueryContract
	ContractConfig
}

type memberQueryHandler struct {
	msg pb.QueryContract
	ContractConfig
}

func NewContractHandler(cfg ContractConfig) *newContractHandler {
	return &newContractHandler{ContractConfig: cfg}
}

func RedeemHandler(cfg ContractConfig) *redeemHandler {
	return &redeemHandler{ContractConfig: cfg}
}

func QueryHandler(cfg ContractConfig) *queryHandler {
	return &queryHandler{ContractConfig: cfg}
}
func MemberQueryHandler(cfg ContractConfig) *memberQueryHandler {
	return &memberQueryHandler{ContractConfig: cfg}
}

func (h *newContractHandler) Msg() proto.Message { return &h.msg }
func (h *redeemHandler) Msg() proto.Message      { return &h.msg }
func (h *queryHandler) Msg() proto.Message       { return &h.msg }
func (h *memberQueryHandler) Msg() proto.Message { return &h.msg }

func (h *newContractHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	if h.pk == nil {
		return nil, errors.New("No publickey")
	}

	contract := make(map[string]uint32)
	for _, m := range h.msg.ContractBody {
		addr, err := txutil.NewAddressFromPBMessage(m.Addr)
		if err != nil {
			return nil, err
		}

		contract[addr.ToString()] = m.Weight
	}

	return h.NewTx(stub, parser.GetNounce()).New(contract, h.pk)
}

func (h *redeemHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := &h.msg

	contract, err := txutil.NewAddressFromPBMessage(msg.Contract)
	if err != nil {
		return nil, err
	}

	redeemAddr := h.redeemAddr
	if redeemAddr == nil {
		redeemAddr, err = txutil.NewAddressFromPBMessage(msg.Redeem)
		if err != nil {
			return nil, err
		}
	}

	var redeemTo []byte

	if to, err := txutil.NewAddressFromPBMessage(msg.To); err == nil {
		redeemTo = to.Hash
	}

	return h.NewTx(stub, parser.GetNounce()).Redeem(contract.Hash, redeemAddr.Hash, toAmount(msg.Amount), redeemTo)
}

func (h *queryHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	msg := &h.msg

	contAddr, err := txutil.NewAddressFromPBMessage(msg.ContractAddr)
	if err != nil {
		return nil, err
	}

	err, data := h.NewTx(stub, parser.GetNounce()).Query(contAddr.Hash)
	if err != nil {
		return nil, err
	}

	return rpc.EncodeRPCResult(data)

}

func (h *memberQueryHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := &h.msg

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

	return rpc.EncodeRPCResult(data)
}
