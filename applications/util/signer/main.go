package main

import (
	"fmt"

	"hyperledger.abchain.org/core/crypto"
	_ "hyperledger.abchain.org/core/crypto/ecdsa"
)

func main() {

	var err error

	var txhashLine, keyLine string

	fmt.Println("Input HASH for tx:")
	_, err = fmt.Scanln(&txhashLine)
	if err != nil {
		panic(err)
	}

	var hash []byte
	_, err = fmt.Sscanf(txhashLine, "%X", &hash)
	if err != nil {
		panic(err)
	}

	if len(hash) == 0 {
		panic("empty hash")
	}

	fmt.Println("Import your key:")
	_, err = fmt.Scanln(&keyLine)
	if err != nil {
		panic(err)
	}

	if priv, err := crypto.PrivatekeyFromString(keyLine); err != nil {
		panic(fmt.Errorf("parse privkey [%s] fail: %s", keyLine, err))
	} else if sig, err := priv.Sign(hash); err != nil {
		panic(fmt.Errorf("signing for hash [%X] fail: %s", hash, err))
	} else {
		fmt.Printf("Signing for hash [%X]:\n", hash)

		sigstr, err := crypto.EncodeCompactSignature(sig)
		if err != nil {
			panic(err)
		}

		fmt.Println(sigstr)
	}

}
