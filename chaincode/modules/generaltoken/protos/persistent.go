package ccprotos

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"hyperledger.abchain.org/core/utils"
)

type NonceKey []byte

func NonceKeyFromString(ns string) (NonceKey, error) {
	var v []byte
	if _, err := fmt.Sscanf(ns, "%X", &v); err != nil {
		return nil, err
	}
	return NonceKey(v), nil
}

func (n NonceKey) MarshalJSON() ([]byte, error) {
	if len(n) == 0 {
		return json.Marshal(nil)
	} else {
		return json.Marshal(fmt.Sprintf("%X", []byte(n)))
	}
}

type FuncRecord_s struct {
	Noncekey []byte
	IsSend   bool
}

func (n *FuncRecord_s) MarshalJSON() ([]byte, error) {
	if len(n.Noncekey) == 0 {
		return json.Marshal(nil)
	}

	if n.IsSend {
		return json.Marshal(fmt.Sprintf("-> %X", n.Noncekey))
	} else {
		return json.Marshal(fmt.Sprintf("<- %X", n.Noncekey))
	}
}

func (n *FuncRecord_s) LoadFromPB(p *FuncRecord) {
	n.IsSend = p.GetIsSend()
	n.Noncekey = p.GetNoncekey()
}

func (n *FuncRecord_s) ToPB() *FuncRecord {

	if n == nil {
		return &FuncRecord{}
	}

	return &FuncRecord{
		Noncekey: n.Noncekey,
		IsSend:   n.IsSend,
	}
}

type accountData_Store struct {
	Balance  *big.Int     `json:"balance"`
	LastFund FuncRecord_s `json:"lastFundID"`
}

type nonceData_Store struct {
	Txid      string       `asn1:"printable" json:"txID"`
	Amount    *big.Int     `json:"amount"`
	FromLast  FuncRecord_s `json:"from"`
	ToLast    FuncRecord_s `json:"to"`
	NonceTime time.Time    `asn1:"generalized" json:"txTime"`
}

type tokenGlobalData_Store struct {
	TotalTokens      *big.Int `json:"total"`
	UnassignedTokens *big.Int `json:"unassign"`
}

type NonceData_s struct {
	nonceData_Store
}

func (n *NonceData_s) GetObject() interface{} { return &n.nonceData_Store }
func (n *NonceData_s) Load(interface{}) error { return nil }
func (n *NonceData_s) Save() interface{}      { return n.nonceData_Store }

func (n *NonceData_s) LoadFromPB(p *NonceData) {

	n.Txid = p.GetTxid()
	n.Amount = big.NewInt(0).SetBytes(p.GetAmount())
	n.FromLast.LoadFromPB(p.GetFromLast())
	n.ToLast.LoadFromPB(p.GetToLast())
	n.NonceTime = utils.ConvertPBTimestamp(p.GetNonceTime())
}

func (n *NonceData_s) ToPB() *NonceData {

	if n == nil {
		return &NonceData{}
	}

	return &NonceData{
		Txid:     n.Txid,
		Amount:   n.Amount.Bytes(),
		FromLast: n.FromLast.ToPB(),
		ToLast:   n.ToLast.ToPB(),
		Other:    &NonceData_NonceTime{NonceTime: utils.CreatePBTimestamp(n.NonceTime)},
	}
}

type AccountData_s struct {
	accountData_Store
}

func (n *AccountData_s) GetObject() interface{} { return &n.accountData_Store }
func (n *AccountData_s) Load(interface{}) error { return nil }
func (n *AccountData_s) Save() interface{}      { return n.accountData_Store }

func (n *AccountData_s) LoadFromPB(p *AccountData) {
	n.Balance = big.NewInt(0).SetBytes(p.GetBalance())
	n.LastFund.LoadFromPB(p.GetLastFund())
}

func (n *AccountData_s) ToPB() *AccountData {

	if n == nil {
		return &AccountData{}
	}

	return &AccountData{
		Balance:  n.Balance.Bytes(),
		LastFund: n.LastFund.ToPB(),
	}
}

type TokenGlobalData_s struct {
	tokenGlobalData_Store
}

func (n *TokenGlobalData_s) GetObject() interface{} { return &n.tokenGlobalData_Store }
func (n *TokenGlobalData_s) Load(interface{}) error { return nil }
func (n *TokenGlobalData_s) Save() interface{}      { return n.tokenGlobalData_Store }

func (n *TokenGlobalData_s) LoadFromPB(p *TokenGlobalData) {
	n.TotalTokens = big.NewInt(0).SetBytes(p.GetTotalTokens())
	n.UnassignedTokens = big.NewInt(0).SetBytes(p.GetUnassignedTokens())
}

func (n *TokenGlobalData_s) ToPB() *TokenGlobalData {

	if n == nil {
		return &TokenGlobalData{}
	}
	return &TokenGlobalData{
		TotalTokens:      n.TotalTokens.Bytes(),
		UnassignedTokens: n.UnassignedTokens.Bytes(),
	}
}
