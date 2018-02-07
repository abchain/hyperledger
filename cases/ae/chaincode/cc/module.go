package chaincode

import (
	"fmt"

	"github.com/abchain/fabric/core/chaincode/shim"
	log "github.com/abchain/fabric/peerex/logging"
	token "hyperledger.abchain.org/chaincode/generaltoken"
	"hyperledger.abchain.org/chaincode/generaltoken/nonce"
	"hyperledger.abchain.org/chaincode/lib/caller"
	"hyperledger.abchain.org/chaincode/lib/txhandle"
	reg "hyperledger.abchain.org/chaincode/registrar"
	share "hyperledger.abchain.org/chaincode/sharesubscription"
)

var logger = log.InitLogger("chaincode")

type AECC struct {
	DebugMode                     bool
	updated                       bool
	assignPrivilege, regPrivilege string
}

const (
	CC_NAME = "AtomicEnergy_v1"
	CC_TAG  = "AE"

	PrivilegeAttr = "Role"
	RegionAttr    = "Region"
)

var invokeMapper map[string]*tx.ChaincodeTx
var queryMapper map[string]*tx.ChaincodeTx
var noncecfg = &nonce.StandardNonceConfig{CC_TAG, false}
var nonceQuerycfg = &nonce.StandardNonceConfig{CC_TAG, true}
var tokencfg = &token.StandardTokenConfig{*noncecfg}
var tokenQuerycfg = &token.StandardTokenConfig{*nonceQuerycfg}
var registrarcfg = &reg.StandardRegistrarConfig{CC_TAG, false, PrivilegeAttr, RegionAttr}
var registrarQuerycfg = &reg.StandardRegistrarConfig{CC_TAG, true, PrivilegeAttr, RegionAttr}
var sharecfg = &share.StandardContractConfig{CC_TAG, false, tokencfg}
var shareQuerycfg = &share.StandardContractConfig{CC_TAG, true, tokenQuerycfg}

func init() {

	invokeMapper = make(map[string]*tx.ChaincodeTx)
	queryMapper = make(map[string]*tx.ChaincodeTx)

	//fundTx hander
	fundH := token.TransferHandler(tokencfg)
	fundTx := &tx.ChaincodeTx{CC_NAME, fundH, nil, nil}
	//transferHandler is also an prehandler
	fundTx.PreHandlers = append(fundTx.PreHandlers, fundH)
	//publickey reg policy
	fundTx.PreHandlers = append(fundTx.PreHandlers,
		reg.RegistrarPreHandler(registrarQuerycfg, fundH))
	invokeMapper[token.Method_Transfer] = fundTx

	assignTx := &tx.ChaincodeTx{CC_NAME, token.AssignHandler(tokencfg), nil, nil}
	invokeMapper[token.Method_Assign] = assignTx

	regTx := &tx.ChaincodeTx{CC_NAME, reg.AdminRegistrarHandler(registrarcfg), nil, nil}
	invokeMapper[reg.Method_AdminRegistrar] = regTx

	newContractH := share.NewContractHandler(sharecfg)
	newContractTx := &tx.ChaincodeTx{CC_NAME, newContractH, nil, nil}
	newContractTx.PreHandlers = append(newContractTx.PreHandlers,
		reg.RegistrarPreHandler(registrarQuerycfg, newContractH))
	newContractTx.PreHandlers = append(newContractTx.PreHandlers, newContractH)
	invokeMapper[share.Method_NewContract] = newContractTx

	redeemH := share.RedeemHandler(sharecfg)
	redeemTx := &tx.ChaincodeTx{CC_NAME, redeemH, nil, nil}
	// Never restrict the redeem address is registred
	// redeemTx.PreHandlers = append(redeemTx.PreHandlers,
	// 	reg.RegistrarPreHandler(registrarQuerycfg, redeemH))
	redeemTx.PreHandlers = append(redeemTx.PreHandlers, redeemH)
	invokeMapper[share.Method_Redeem] = redeemTx

	queryMapper[token.Method_QueryToken] = &tx.ChaincodeTx{CC_NAME, token.TokenQueryHandler(tokenQuerycfg), nil, nil}
	queryMapper[token.Method_QueryGlobal] = &tx.ChaincodeTx{CC_NAME, token.GlobalQueryHandler(tokenQuerycfg), nil, nil}
	queryMapper[token.Method_QueryTrans] = &tx.ChaincodeTx{CC_NAME, nonce.NonceQueryHandler(nonceQuerycfg), nil, nil}
	queryMapper[reg.Method_Registrar] = &tx.ChaincodeTx{CC_NAME, reg.QueryPkHandler(registrarQuerycfg), nil, nil}
	//TODO: add suitable policy for contract query
	queryMapper[share.Method_Query] = &tx.ChaincodeTx{CC_NAME, share.QueryHandler(shareQuerycfg), nil, nil}
	queryMapper[share.Method_MemberQuery] = &tx.ChaincodeTx{CC_NAME, share.MemberQueryHandler(shareQuerycfg), nil, nil}

}

func (t *AECC) updateGlobal(stub shim.ChaincodeStubInterface) error {
	if t.DebugMode || t.updated {
		return nil
	}

	err, data := registrarQuerycfg.NewTx(stub).Global()
	if err != nil {
		return err
	}

	t.assignPrivilege = data.AdminPrivilege
	t.regPrivilege = data.RegPrivilege

	//update related handler
	h, ok := invokeMapper[token.Method_Assign]
	if ok {
		var assignPriv tx.TxAttrVerifier = make(map[string]string)
		assignPriv[PrivilegeAttr] = t.assignPrivilege
		h.PreHandlers = append(h.PreHandlers, assignPriv)
	}

	t.updated = true
	return nil
}

func (t *AECC) Init(stub shim.ChaincodeStubInterface,
	function string, args []string) ([]byte, error) {

	if function == "INIT" {
		handlers := make(map[string]rpc.DeployHandler)

		handlers[token.DeployMethod] = token.CCDeployHandler(CC_TAG)
		handlers[reg.DeployMethod] = reg.CCDeployHandler(CC_TAG)

		err := rpc.DeployCC(stub, args, handlers)
		if err != nil {
			return nil, err
		}

		return []byte("OK"), nil
	} else {
		return nil, fmt.Errorf("Not support method %s", function)
	}

}

func (t *AECC) Invoke(stub shim.ChaincodeStubInterface,
	function string, args []string) ([]byte, error) {

	err := t.updateGlobal(stub)
	if err != nil {
		return nil, err
	}

	h, ok := invokeMapper[function]
	if !ok {
		return nil, fmt.Errorf("Method not found for method %s", function)
	}

	return h.TxCall(stub, function, args)
}

func (t *AECC) Query(stub shim.ChaincodeStubInterface,
	function string, args []string) ([]byte, error) {

	err := t.updateGlobal(stub)
	if err != nil {
		return nil, err
	}

	h, ok := queryMapper[function]
	if !ok {
		return nil, fmt.Errorf("Method not found for method %s", function)
	}

	return h.TxCall(stub, function, args)
}

func ExportMain() {
	err := shim.Start(new(AECC))
	if err != nil {
		fmt.Printf("Error starting AE chaincode: %s", err)
	}
}
