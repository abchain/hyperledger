package ccprotos

import (
	"hyperledger.abchain.org/core/utils"
	"math/big"
	"time"
)

type FuncRecord_s struct {
	Noncekey []byte
	IsSend   bool
}

func (n *FuncRecord_s) LoadFromPB(p *FuncRecord) {
	n.IsSend = p.GetIsSend()
	n.Noncekey = p.GetNoncekey()
}

func (n *FuncRecord_s) ToPB() *FuncRecord {

	if n == nil {
		return nil
	}

	return &FuncRecord{
		Noncekey: n.Noncekey,
		IsSend:   n.IsSend,
	}
}

type accountData_Store struct {
	Balance  *big.Int
	LastFund FuncRecord_s
}

type nonceData_Store struct {
	Txid      string
	Amount    *big.Int
	FromLast  FuncRecord_s
	ToLast    FuncRecord_s
	NonceTime time.Time
}

type tokenGlobalData_Store struct {
	TotalTokens      *big.Int
	UnassignedTokens *big.Int
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
		return nil
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
		return nil
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
		return nil
	}

	return &TokenGlobalData{
		TotalTokens:      n.TotalTokens.Bytes(),
		UnassignedTokens: n.UnassignedTokens.Bytes(),
	}
}
