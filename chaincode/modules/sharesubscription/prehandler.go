package subscription

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/lib/runtime"
	"hyperledger.abchain.org/chaincode/lib/txhandle"
	pb "hyperledger.abchain.org/chaincode/modules/sharesubscription/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
)

type redeemContractAddrCred struct {
	ContractConfig
}

func NewRedeemContractAddrCred(cfg ContractConfig) redeemContractAddrCred {
	return redeemContractAddrCred{cfg}
}

//this module will first catch the runtime data (as prehandler) to collect redeem address (if required)
func (v redeemContractAddrCred) PreHandling(stub shim.ChaincodeStubInterface, _ string, parser txutil.Parser) error {

	m, ok := parser.GetMessage().(*pb.RedeemContract)
	if !ok {
		return errors.New("Binding to wrong txhandler")
	}

	rt := v.NewTx(stub, parser.GetNonce())
	err, ct := rt.Query_C(m.GetContract().GetHash())
	if err != nil {
		return err
	}

	if len(m.GetRedeems()) == 0 {

		cred := parser.GetAddrCredential()
		if cred == nil {
			return errors.New("Could not found which addresses should be redeem to")
		}

		for _, pk := range cred.ListCredPubkeys() {
			addr, err := txutil.NewAddress(pk)
			if err != nil {
				return err
			}

			if _, ok := ct.Find(addr.Internal()); ok {
				m.Redeems = append(m.Redeems, addr.PBMessage())
			}
		}

	} else {
		//match
		for _, addr := range m.GetRedeems() {
			caddr, err := txutil.NewAddressFromPBMessage(addr)
			if err == nil {
				if _, ok := ct.Find(caddr.Internal()); !ok {
					return errors.New("Invalid redeem addr (not in contract)")
				}
			} else {
				return err
			}
		}
	}

	return nil
}

//notice this can not be applied on external invoking
type contractCred struct {
	*StandardContractConfig
}

func NewContractVerifier(cfg *StandardContractConfig) contractCred {
	return contractCred{cfg}
}

func (v contractCred) PostHandling(stub shim.ChaincodeStubInterface, _ string,
	parser txutil.Parser, ret []byte) ([]byte, error) {

	rt := runtime.NewRuntime(v.Root, stub, v.Config).SubRuntime(contract_auth_tag)
	msg := parser.GetMessage().(*pb.RegContract)

	addrhash, err := v.AddrCfg.NewTx(stub, parser.GetNonce()).NormalizeAddress(ret)
	if err != nil {
		return nil, err
	}

	degaddr, err := txutil.NewAddressFromPBMessage(msg.GetDelegator())
	if err != nil {
		return nil, err
	}

	err = rt.Storage.SetRaw(addrToKey(addrhash), degaddr.Internal())
	if err != nil {
		return nil, err
	}

	return ret, nil
}

//an address interface in verifier, checking for the contract's deligator
func GetDeletagorAddress(cfg *StandardContractConfig) func(shim.ChaincodeStubInterface,
	proto.Message) []*txutil.Address {

	return func(stub shim.ChaincodeStubInterface,
		msg proto.Message) (ret []*txutil.Address) {

		if ma, ok := msg.(tx.MsgAddresses); !ok {
			return
		} else {
			ret = ma.GetAddresses()
		}

		rt := runtime.NewRuntime(cfg.Root, stub, cfg.Config).SubRuntime(contract_auth_tag)

		for i, addr := range ret {
			deletagorAddr, err := rt.Storage.GetRaw(addrToKey(addr.Internal()))
			if err == nil && len(deletagorAddr) != 0 {
				ret[i] = txutil.NewAddressFromHash(deletagorAddr)
			}
		}
		return
	}
}
