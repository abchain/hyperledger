package main

import (
	"hyperledger.abchain.org/client"
	_ "hyperledger.abchain.org/client/yafabric"
	"hyperledger.abchain.org/core/config"

	"encoding/json"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"math/big"
	"os"
	"strings"
)

var configTitle string
var blockspec string
var ccFilter string
var includeErrTx bool
var progSpec int

type argsDump struct {
	Function string   `json:"function,omitempty"`
	Args     [][]byte `json:"args,omitempty"`
	ByteArgs bool
}

func main() {

	flag.StringVar(&configTitle, "config", "config", "config file name")
	flag.StringVar(&blockspec, "block", "", "range of block: [n] or [start-end]")
	flag.StringVar(&ccFilter, "chaincodes", "", "filter by chaincode name, separated by comma")
	flag.BoolVar(&includeErrTx, "dumpErrTx", false, "also dump the tx which is errored")
	flag.IntVar(&progSpec, "printprog", -1, "print a dot after each N blocks is dumped")
	flag.Parse()

	if err := config.LoadConfig(configTitle, nil); err != nil {
		panic(err)
	}

	cfg := client.NewFabricRPCConfig("NOTUSED")
	err := cfg.UseYAFabricCli(viper.GetViper())
	if err != nil {
		panic(err)
	}

	var beg, end int64
	if blockspec != "" {
		val := strings.Split(blockspec, "-")
		if len(val) >= 1 {
			if vn, done := big.NewInt(0).SetString(val[0], 0); !done {
				panic("Wrong blocknumber:" + blockspec)
			} else if beg = int64(vn.Uint64()); beg < 0 {
				panic("Number too large:" + blockspec)
			}

			end = beg

			if len(val) >= 2 {
				if val[1] != "" {
					if vn, done := big.NewInt(0).SetString(val[1], 0); !done {
						panic("Wrong blocknumber:" + blockspec)
					} else if end = int64(vn.Uint64()); end < 0 {
						panic("Number too large:" + blockspec)
					}

					if beg > end {
						panic("Wrong range:" + blockspec)
					}
				} else {
					end = 0
				}

			}
		}
	}

	var ccFilters []string
	if ccFilter != "" {
		ccFilters = strings.Split(ccFilter, ",")
	}

	chain, err := cfg.GetChain()
	if err != nil {
		panic(err)
	}

	if cn, err := chain.GetChain(); err != nil {
		panic(err)
	} else {

		if end == 0 || cn.Height <= end {
			end = cn.Height - 1
		}

		fmt.Fprintf(os.Stderr, "Chain has %d blocks, dumping from %d to %d\n", cn.Height, beg, end)
	}

	var txcnt, failcnt, progcnt int

	for i := beg; i <= end; i++ {
		blk, err := chain.GetBlock(i)
		if err != nil {
			panic(err)
		}

		//construct err tx mapper
		errTxs := map[string]bool{}
		if !includeErrTx {
			for _, evt := range blk.TxEvents {

				if evt.Status > 0 {
					errTxs[evt.TxID] = true
				}

			}
		}

		for _, tx := range blk.Transactions {

			if _, existed := errTxs[tx.TxID]; existed {
				continue
			}

			if len(ccFilters) > 0 {
				passed := false
				for _, ccn := range ccFilters {
					if ccn == tx.Chaincode {
						passed = true
					}
				}

				if !passed {
					continue
				}
			}

			txcnt++
			dumpstr, err := json.Marshal(&argsDump{tx.Method, tx.TxArgs, true})
			if err != nil {
				fmt.Fprintf(os.Stderr, "dumping tx [%v] fail: %s\n", tx, err)
				failcnt++
			} else {
				fmt.Println(string(dumpstr))
			}
		}

		progcnt++
		if progSpec > 0 && progcnt >= progSpec {
			fmt.Fprintf(os.Stderr, ".")
			progcnt = 0
		}
	}

	fmt.Fprintf(os.Stderr, "\nDump %d transactions, %d fail\n", txcnt-failcnt, failcnt)
}
