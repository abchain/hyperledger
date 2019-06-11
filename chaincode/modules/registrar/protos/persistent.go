package ccprotos

import (
	proto "github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/core/crypto"
	"hyperledger.abchain.org/core/utils"
	pb "hyperledger.abchain.org/protos"
	"time"
)

type globalData_Store struct {
	RegPrivilege   string `asn1:"printable"`
	AdminPrivilege string `asn1:"printable"`
	DeployFlag     []byte
}

type RegGlobalData_s struct {
	globalData_Store
}

func (n *RegGlobalData_s) GetObject() interface{} { return &n.globalData_Store }
func (n *RegGlobalData_s) Load(interface{}) error { return nil }
func (n *RegGlobalData_s) Save() interface{}      { return n.globalData_Store }
func (n *RegGlobalData_s) ToPB() *RegGlobalData {
	return &RegGlobalData{
		RegPrivilege:   n.RegPrivilege,
		AdminPrivilege: n.AdminPrivilege,
		DeployFlag:     n.DeployFlag,
	}
}

type regData_Store struct {
	PkBytes   []byte
	RegTxid   string `asn1:"printable"`
	Region    string `asn1:"utf8"`
	Enabled   bool
	RegTs     time.Time `asn1:"generalized"`
	AuthCodes []int32
}

type RegData_s struct {
	regData_Store
	Pk crypto.Verifier
}

func (n *RegData_s) GetObject() interface{} { return &n.regData_Store }
func (n *RegData_s) Save() interface{}      { return n.regData_Store }
func (n *RegData_s) Load(interface{}) error {

	pkpb := new(pb.PublicKey)
	var err error
	err = proto.Unmarshal(n.PkBytes, pkpb)
	if err != nil {
		return err
	}

	n.Pk, err = crypto.PublicKeyFromPBMessage(pkpb)
	if err != nil {
		return err
	}

	return nil
}

func (n *RegData_s) ToPB() *RegData {
	return &RegData{
		Pk:        n.Pk.PBMessage().(*pb.PublicKey),
		RegTxid:   n.RegTxid,
		Region:    n.Region,
		Enabled:   n.Enabled,
		Authcodes: n.AuthCodes,
		RegTs:     utils.CreatePBTimestamp(n.RegTs),
	}
}
