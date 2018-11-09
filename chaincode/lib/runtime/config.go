package runtime

import (
	_ "hyperledger.abchain.org/chaincode/shim"
)

type Config struct {
	ReadOnly bool
}

func NewConfig() *Config { return &Config{} }

func (c *Config) SetReadOnly(flag bool) { c.ReadOnly = flag }
