package txrepair

import (
	"bytes"
	"golang.org/x/net/context"
	"time"

	ccplatform "github.com/abchain/fabric/core/chaincode"
	"github.com/abchain/fabric/core/config"
	"github.com/abchain/fabric/core/db"
	"github.com/abchain/fabric/core/embedded_chaincode/api"
	"github.com/abchain/fabric/core/ledger"
	"github.com/abchain/fabric/core/ledger/genesis"
	"github.com/abchain/fabric/node"
	"github.com/abchain/fabric/node/start"
	"github.com/abchain/fabric/protos"

	"hyperledger.abchain.org/chaincode/impl/yafabric"
	"hyperledger.abchain.org/chaincode/shim"
)

var logger = config.GetLogger("repair")
var CustomGenesisKV = map[string][]byte{}
var CustomGenesisCCName = "unknown"

func Run(configTag string, ccs map[string]shim.Chaincode) {

	gcfg := new(startnode.GlobalConfig)
	gcfg.EnvPrefix = configTag
	gcfg.ConfigFileName = configTag
	if err := gcfg.Apply(); err != nil {
		logger.Fatal(err)
	}

	// Init the crypto layer
	if err := config.InitCryptoGlobal(nil); err != nil {
		logger.Fatalf("Failed to initialize the crypto layer: %s", err)
	}

	ne := new(node.NodeEngine)
	ne.Name = "repairNode"
	defer ne.FinalRelease()
	if err := ne.Init(); err != nil {
		logger.Fatal(err)
	}

	ccplatform.NewSystemChaincodeSupport(ne.Name)
	ccplatform.NewSystemChaincodeSupport(ne.Name, ccplatform.DefaultChain)

	for name, bcc := range ccs {
		cc := fabric_impl.GenYAfabricCC(bcc)
		if err := api.RegisterECC(&api.EmbeddedChaincode{name, cc}); err != nil {
			logger.Fatalf("launch cc fail: %s", err)
		}
		logger.Debugf("Launch chaincode <%s>", name)
	}

	l := ne.DefaultLedger()
	if l.GetBlockchainSize() <= 1 {
		logger.Fatal("Empty ledger, nothing to do")
	}

	worksn, _ := l.CreateSnapshot()
	if worksn == nil {
		logger.Fatal("Empty snapshot, something wrong")
	}
	defer worksn.Release()

	till, err := worksn.TestContinuouslBlockRange()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("Start fix for %d blocks", till)

	genesisblk, err := worksn.GetBlockByNumber(0)
	if err != nil {
		logger.Fatalf("Fail at obtaining genesis block", err)
	}

	err = db.GetDBHandle().DeleteAll()
	if err != nil {
		logger.Fatal(err)
	}

	workLedger, err := ledger.GetNewLedger(db.GetDBHandle(), nil)
	if err != nil {
		logger.Fatal(err)
	}

	err = genesis.NewGenesisBlock(genesisblk).MakeGenesisForLedger(workLedger,
		CustomGenesisCCName, CustomGenesisKV)

	if err != nil {
		logger.Fatal(err)
	}

	parentstate := genesisblk.GetStateHash()
	logger.Infof("Set genesis to be %X", parentstate)

	blkcommiter := ledger.NewBlockAgent(workLedger)

	for i := uint64(1); i <= till; i++ {

		refblk, err := worksn.GetBlockByNumber(i)
		if err != nil {
			logger.Fatalf("Fail at block %d: %s", i, err)
		}
		oldstate := refblk.GetStateHash()

		var txes []*protos.TransactionHandlingContext
		for _, tx := range refblk.GetTransactions() {
			txe, _ := protos.DefaultTxHandler.Handle(protos.NewTransactionHandlingContext(tx))
			txes = append(txes, txe)
		}

		txagent, _ := ledger.NewTxEvaluatingAgent(workLedger)

		blockTs := time.Unix(0, 0)
		if blkts := refblk.GetTimestamp(); blkts != nil {
			blockTs = protos.GetUnixTime(blkts)
		}

		if _, err := ccplatform.ExecuteTransactions2(context.Background(), ccplatform.DefaultChain,
			txes, blockTs, txagent); err != nil {
			logger.Fatalf("Fail at exec txs on %d: %s", i, err)
		}

		refblk, err = txagent.PreviewBlock(i, refblk)
		if err != nil {
			logger.Fatalf("Fail at preview block %d: %s", i, err)
		}

		if csh := refblk.GetStateHash(); bytes.Compare(oldstate, csh) != 0 {
			logger.Warningf("state on %d is change (%.16X to %.16X)", i, oldstate, csh)
			err = workLedger.AddGlobalState(parentstate, csh)
			if err != nil {
				logger.Errorf("add global state fail: %s", err)
			}
			parentstate = csh
		} else {
			logger.Infof("state on %d is identical (%.16X)", i, csh)
			parentstate = csh
		}

		err = blkcommiter.SyncCommitBlock(i, refblk)
		if err != nil {
			logger.Fatalf("Fail at commit block %d: %s", i, err)
		}

		err = txagent.StateCommitOne(i, refblk)
		if err != nil {
			logger.Fatalf("Fail at commit state %d: %s", i, err)
		}
	}

}
