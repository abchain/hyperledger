package crypto

import (
	"fmt"
	"hyperledger.abchain.org/protos"
	"math/big"
	"strings"
)

// <curvetype>,<pub part>,<sig part>
func decodeECDSASig(sig string) (*protos.Signature_ECDSA, error) {
	sigparts := strings.Split(sig, ",")
	if len(sigparts) != 3 {
		return nil, fmt.Errorf("Invalid encoding for ecdsa sig")
	}

	ret := &protos.Signature_ECDSA{}
	_, err := fmt.Sscanf(sigparts[0], "%d", &ret.Curvetype)
	if err != nil {
		return nil, fmt.Errorf("Parse curve type fail")
	}

	if len(sigparts[1]) < 3 {
		//pubk part is considered as V
		v := &protos.Signature_ECDSA_V{}
		_, err = fmt.Sscanf(sigparts[1], "%d", &v.V)
		if err != nil {
			return nil, fmt.Errorf("Parse sig.v fail")
		}

		ret.Pub = v
	} else {
		//pubk part is considered as ECPoint
		var bt []byte
		_, err = fmt.Sscanf(sigparts[1], "%x", &bt)
		if err != nil {
			return nil, fmt.Errorf("Parse sig.p bytes fail")
		} else if len(bt) != 64 {
			return nil, fmt.Errorf("Parse sig.p bytes has no expected length [%d]", len(bt))
		}

		p := &protos.Signature_ECDSA_P{P: &protos.ECPoint{X: bt[:32], Y: bt[32:]}}
		ret.Pub = p
	}

	var bt []byte
	_, err = fmt.Sscanf(sigparts[2], "%x", &bt)
	if err != nil {
		return nil, fmt.Errorf("Parse sig bytes fail")
	} else if len(bt) != 64 {
		return nil, fmt.Errorf("Parse sig bytes has no expected length [%d]", len(bt))
	}

	ret.R = bt[:32]
	ret.S = bt[32:]

	return ret, nil
}

func EncodeCompactSignature(sig *protos.Signature) (string, error) {

	var scheme, sigdata, derv string

	switch ssig := sig.Data.(type) {
	case *protos.Signature_Ec:
		scheme = "EC"
		ecsig := ssig.Ec
		var pkorv string
		if ecp := ecsig.GetP(); ecp != nil {
			pkorv = fmt.Sprintf("%X%X", ecp.GetX(), ecp.GetY())
		} else {
			pkorv = fmt.Sprintf("%d", ecsig.GetV())
		}
		sigdata = fmt.Sprintf("%02d,%s,%X%X", ecsig.GetCurvetype(), pkorv, ecsig.GetR(), ecsig.GetS())
	default:
		return "", fmt.Errorf("No signature data")
	}

	if kd := sig.Kd; kd != nil {
		indstr := big.NewInt(0).SetBytes(kd.GetIndex()).String()
		if cc := kd.GetChaincode(); cc != nil {
			derv = fmt.Sprintf("%s,%X,%X", indstr, kd.GetRootFingerprint(), cc)
		} else {
			derv = fmt.Sprintf("%s,%X", indstr, kd.GetRootFingerprint())
		}
	}

	return strings.Join([]string{scheme, sigdata, derv}, ":"), nil
}

//we define a compact, simple encoding for the protbuf signature before our SDK for other language is mature...
//which has such a form: <crypto scheme mark>:<encoded sig>:[derived part := <index>,<finger print>,[cc]]
func DecodeCompactSignature(sig string) (*protos.Signature, error) {

	sigparts := strings.Split(sig, ":")
	if len(sigparts) != 3 {
		return nil, fmt.Errorf("Invalid encoding for sig")
	}

	ret := &protos.Signature{}

	switch sigparts[0] {
	case "EC":
		sigec, err := decodeECDSASig(sigparts[1])
		if err != nil {
			return nil, fmt.Errorf("Decode ecdsa sig fail: %s", err)
		}

		ret.Data = &protos.Signature_Ec{sigec}
	default:
		return nil, fmt.Errorf("Unknown sig scheme: %s", sigparts[0])
	}

	if len(sigparts[2]) != 0 {
		dervParts := strings.Split(sigparts[2], ",")
		if len(dervParts) < 2 {
			return nil, fmt.Errorf("Invalid encoding for key derived data")
		}
		ret.Kd = new(protos.KeyDerived)

		if ind, ok := big.NewInt(0).SetString(dervParts[0], 10); !ok {
			return nil, fmt.Errorf("Decode key derived index fail")
		} else {
			ret.Kd.Index = ind.Bytes()
		}

		if _, err := fmt.Sscanf(dervParts[1], "%x", &ret.Kd.RootFingerprint); err != nil {
			return nil, fmt.Errorf("Decode key derived root fp fail: %s", err)
		}

		if len(dervParts) > 2 {
			if _, err := fmt.Sscanf(dervParts[2], "%x", &ret.Kd.Chaincode); err != nil {
				return nil, fmt.Errorf("Decode key derived root fp fail: %s", err)
			}
		}
	}

	return ret, nil
}
