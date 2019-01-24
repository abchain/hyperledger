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

type registrarHandler struct{ RegistrarConfig }
type adminRegistrarHandler struct{ RegistrarConfig }
type revokePkHandler struct{ RegistrarConfig }
type activePkHandler struct{ RegistrarConfig }
type queryPkHandler struct{ RegistrarConfig }
type initHandler struct{ RegistrarConfig }

func RegistrarHandler(cfg RegistrarConfig) registrarHandler {
	return registrarHandler{RegistrarConfig: cfg}
}

func AdminRegistrarHandler(cfg RegistrarConfig) adminRegistrarHandler {
	return adminRegistrarHandler{RegistrarConfig: cfg}
}

func RevokePkHandler(cfg RegistrarConfig) revokePkHandler {
	return revokePkHandler{RegistrarConfig: cfg}
}
func ActivePkHandler(cfg RegistrarConfig) activePkHandler {
	return activePkHandler{RegistrarConfig: cfg}
}

func QueryPkHandler(cfg RegistrarConfig) queryPkHandler {
	return queryPkHandler{RegistrarConfig: cfg}
}

func InitHandler(cfg RegistrarConfig) initHandler {
	return initHandler{RegistrarConfig: cfg}
}

func (h registrarHandler) Msg() proto.Message      { return new(ccpb.RegPublicKey) }
func (h adminRegistrarHandler) Msg() proto.Message { return new(ccpb.RegPublicKey) }
func (h revokePkHandler) Msg() proto.Message       { return new(ccpb.RevokePublicKey) }
func (h activePkHandler) Msg() proto.Message       { return new(ccpb.ActivePublicKey) }
func (h queryPkHandler) Msg() proto.Message        { return new(ccpb.ActivePublicKey) }
func (h initHandler) Msg() proto.Message           { return new(ccpb.Settings) }

func (h registrarHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := parser.GetMessage().(*ccpb.RegPublicKey)
	err := h.NewTx(stub).Registrar(msg.PkBytes, msg.Region)
	if err != nil {
		return nil, err
	}

	return []byte("OK"), nil
}

func (h adminRegistrarHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := parser.GetMessage().(*ccpb.RegPublicKey)

	err := h.NewTx(stub).AdminRegistrar(msg.PkBytes)
	if err != nil {
		return nil, err
	}
	return []byte("OK"), nil
}

func (h revokePkHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := parser.GetMessage().(*ccpb.RevokePublicKey)

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

func (h activePkHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	msg := parser.GetMessage().(*ccpb.ActivePublicKey)

	err := h.NewTx(stub).ActivePk(msg.Key)
	if err != nil {
		return nil, err
	}

	return []byte("OK"), nil
}

func (h queryPkHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	msg := parser.GetMessage().(*ccpb.ActivePublicKey)

	err, data := h.NewTx(stub).Pubkey(msg.Key)
	if err != nil {
		return nil, err
	}

	return rpc.EncodeRPCResult(data)
}

func (h initHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := parser.GetMessage().(*ccpb.Settings)

	reg := h.NewTx(stub)

	err := reg.Init(!msg.DebugMode, msg.AdminPrivilege, msg.RegPrivilege)
	if err != nil {
		return nil, err
	}

	return []byte("Ok"), nil
}
