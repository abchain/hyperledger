package client

import (
	"encoding/json"
	"errors"
	"fmt"
	protos "github.com/abchain/fabric/protos"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	"io/ioutil"
	"math/big"
	"net/http"
)

type RESTClient struct {
	server string
}

func NewRESTClient(server string) (*RESTClient, error) {

	c := &RESTClient{server}

	return c, nil
}

type restError struct {
	Error string `json:"Error,omitempty"`
}

func (c *RESTClient) GetBlockchainInfo() (*protos.BlockchainInfo, error) {

	if c.server == "" {
		return nil, errors.New("REST Server is not set")
	}

	// Generate URL
	url := fmt.Sprintf("%s/chain", c.server)
	//logger.Debugf("Request URL(%v)", url)

	// HTTP Request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Requset failed: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read response failed: %v", err)
	}
	//logger.Debugf("Body: %v", string(body))

	// Parse error
	errMsg := &restError{}
	json.Unmarshal(body, errMsg)
	if errMsg.Error != "" {
		return nil, fmt.Errorf("Request failed: %v", errMsg.Error)
	}

	// Unmarshal
	info1 := &protos.BlockchainInfo{}
	err = json.Unmarshal(body, info1)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal response failed: %v", err)
	}

	return info1, nil
}

func (c *RESTClient) GetBlock(height *big.Int) (*protos.Block, error) {

	if c.server == "" {
		return nil, errors.New("REST Server is not set")
	}

	// Generate URL
	url := fmt.Sprintf("%s/chain/blocks/%v", c.server, height.Int64())
	//logger.Debugf("Request URL(%v)", url)

	// HTTP Request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Requset failed: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read response failed: %v", err)
	}
	//logger.Debugf("Body: %v", string(body))

	// Parse error
	errMsg := &restError{}
	json.Unmarshal(body, errMsg)
	if errMsg.Error != "" {
		return nil, fmt.Errorf("Request failed: %v", errMsg.Error)
	}

	// Unmarshal
	block1 := &protos.Block{}
	err = json.Unmarshal(body, block1)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal response failed: %v", err)
	}
	return block1, nil
}

func (c *RESTClient) GetTransaction(transactionID string) (*protos.Transaction, error) {

	if c.server == "" {
		return nil, errors.New("REST Server is not set")
	}

	// Generate URL
	url := fmt.Sprintf("%s/transactions/%v", c.server, transactionID)
	//logger.Debugf("Request URL(%v)", url)

	// HTTP Request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Requset failed: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read response failed: %v", err)
	}
	//logger.Debugf("Body: %v", string(body))

	// Parse error
	errMsg := &restError{}
	json.Unmarshal(body, errMsg)
	if errMsg.Error != "" {
		return nil, fmt.Errorf("Request failed: %v", errMsg.Error)
	}

	// Unmarshal
	tx1 := &protos.Transaction{}
	err = json.Unmarshal(body, tx1)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal response failed: %v", err)
	}

	return tx1, nil
}
