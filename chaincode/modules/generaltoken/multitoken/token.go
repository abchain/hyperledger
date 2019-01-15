package multitoken

import (
	"errors"
	"hyperledger.abchain.org/chaincode/modules/generaltoken"
	"math/big"
	"regexp"
)

var baseVerifier = regexp.MustCompile(`[A-Za-z0-9]{4,}`)

func baseNameVerifier(name string) error {
	ret := baseVerifier.FindString(name)
	if len(ret) < len(name) {
		return errors.New("Token name contain invalid part")
	}

	return nil
}

func (mtoken *baseMultiTokenTx) GetToken(name string) (generaltoken.TokenTx, error) {

	if err := baseNameVerifier(name); err != nil {
		return nil, err
	}

	subrt := mtoken.ChaincodeRuntime.SubRuntime(name)

	return generaltoken.NewTokenTxImpl(subrt, mtoken.nonce, mtoken.tokenNonce), nil
}

func (mtoken *baseMultiTokenTx) CreateToken(name string, total *big.Int) error {

	if err := baseNameVerifier(name); err != nil {
		return err
	}

	tk := generaltoken.NewTokenTxImpl(mtoken.ChaincodeRuntime.SubRuntime(name), mtoken.nonce, mtoken.tokenNonce)
	return tk.Init(total)
}
