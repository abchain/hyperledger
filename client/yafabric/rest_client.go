package client

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"hyperledger.abchain.org/client"
	"hyperledger.abchain.org/client/yafabric/protos"
	"io/ioutil"
	"net/http"
)

type restcli string

type restError struct {
	Error string `json:"Error,omitempty"`
}

func queryBlockchain(url string, out interface{}) error {
	// HTTP Request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Requset failed: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Read response failed: %v", err)
	}
	//logger.Debugf("Body: %v", string(body))

	// Parse error
	errMsg := &restError{}
	json.Unmarshal(body, errMsg)
	if errMsg.Error != "" {
		return fmt.Errorf("Request failed: %v", errMsg.Error)
	}

	// Unmarshal
	err = json.Unmarshal(body, out)
	if err != nil {
		return fmt.Errorf("Unmarshal response failed: %v", err)
	}

	return nil
}

func (server restcli) GetCurrentBlock() (int64, error) {

	info := &protos.BlockchainInfo{}
	err := queryBlockchain(fmt.Sprintf("http://%s/chain", server), info)
	if err != nil {
		return 0, err
	}
	return int64(info.GetHeight()), nil
}

func (server restcli) GetBlock(h int64) (*protos.Block, error) {

	block := &protos.Block{}
	err := queryBlockchain(fmt.Sprintf("http://%s/chain/blocks/%d", server, h), block)
	if err != nil {
		return nil, err
	}

	return block, nil

}

func (server restcli) GetTransaction(transactionID string) (*protos.Transaction, error) {

	tx := &protos.Transaction{}
	err := queryBlockchain(fmt.Sprintf("http://%s/transactions/%s", server, transactionID), tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (server restcli) GetTxIndex(string) (int64, error) {
	//REST have no implement for index yet
	return 0, nil
}

type restClient struct {
}

func (restClient) ViaWeb(vp *viper.Viper) client.ChainInfo {

	srv := vp.GetString("server")
	if srv == "" {
		srv = "localhost:8080"
	}

	return &blockchainInterpreter{restcli(srv)}
}

func NewRESTConfig() client.ChainClient { return restClient{} }

func init() {
	client.ChainProxy_Impls["yafabric"] = NewRESTConfig
}
