package registrar

import (
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/lib/caller"
	ccpb "hyperledger.abchain.org/chaincode/modules/registrar/protos"
	"hyperledger.abchain.org/chaincode/shim"
	"hyperledger.abchain.org/core/crypto"
	txutil "hyperledger.abchain.org/core/tx"
)

//wrap ccpb.RegPublicKey with ParseAddress interface
type RegPkMsg struct {
	msg ccpb.RegPublicKey
}

type registrarHandler struct {
	RegPkMsg
	RegistrarConfig
}

type adminRegistrarHandler struct {
	RegPkMsg
	RegistrarConfig
}
type revokePkHandler struct {
	msg ccpb.RevokePublicKey
	RegistrarConfig
}

type activePkHandler struct {
	msg ccpb.ActivePublicKey
	RegistrarConfig
}

type queryPkHandler struct {
	msg ccpb.ActivePublicKey
	RegistrarConfig
}

type initHandler struct {
	msg ccpb.Settings
	RegistrarConfig
}

func RegistrarHandler(cfg RegistrarConfig) *registrarHandler {
	return &registrarHandler{RegistrarConfig: cfg}
}

func AdminRegistrarHandler(cfg RegistrarConfig) *adminRegistrarHandler {
	return &adminRegistrarHandler{RegistrarConfig: cfg}
}

func RevokePkHandler(cfg RegistrarConfig) *revokePkHandler {
	return &revokePkHandler{RegistrarConfig: cfg}
}
func ActivePkHandler(cfg RegistrarConfig) *activePkHandler {
	return &activePkHandler{RegistrarConfig: cfg}
}

func QueryPkHandler(cfg RegistrarConfig) *queryPkHandler {
	return &queryPkHandler{RegistrarConfig: cfg}
}

func InitHandler(cfg RegistrarConfig) *initHandler {
	return &initHandler{RegistrarConfig: cfg}
}

func (h *registrarHandler) Msg() proto.Message      { return &h.msg }
func (h *adminRegistrarHandler) Msg() proto.Message { return &h.msg }
func (h *revokePkHandler) Msg() proto.Message       { return &h.msg }
func (h *activePkHandler) Msg() proto.Message       { return &h.msg }
func (h *queryPkHandler) Msg() proto.Message        { return &h.msg }
func (h *initHandler) Msg() proto.Message           { return &h.msg }

func (h *registrarHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := &h.msg
	pk, err := crypto.PublicKeyFromPBMessage(msg.Pk)
	if err != nil {
		return nil, err
	}

	return h.NewTx(stub).Registrar(pk, msg.Region)
}

func (h *adminRegistrarHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := &h.msg

	pk, err := crypto.PublicKeyFromPBMessage(msg.Pk)
	if err != nil {
		return nil, err
	}

	err = h.NewTx(stub).AdminRegistrar(pk)
	if err != nil {
		return nil, err
	}
	return []byte("OK"), nil
}

func (h *revokePkHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := &h.msg

	pk, err := crypto.PublicKeyFromPBMessage(msg.Pk)
	if err != nil {
		return nil, err
	}

	err = h.NewTx(stub).RevokePk(pk)
	if err != nil {
		return nil, err
	}
	return []byte("OK"), nil
}

func (h *activePkHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	msg := &h.msg

	err := h.NewTx(stub).ActivePk(msg.Key)
	if err != nil {
		return nil, err
	}

	return []byte("OK"), nil
}

func (h *queryPkHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	msg := &h.msg

	err, data := h.NewTx(stub).Pubkey(msg.Key)
	if err != nil {
		return nil, err
	}

	return rpc.EncodeRPCResult(data)
}

func (h *initHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := &h.msg

	reg := h.NewTx(stub)

	err := reg.Init(!msg.DebugMode, msg.AdminPrivilege, msg.RegPrivilege)
	if err != nil {
		return nil, err
	}

	return []byte("Ok"), nil
}
