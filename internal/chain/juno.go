package chain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Juno struct {
}

func (Juno) GetTx(txHash string) (*TxResult, error) {
	txHash = "0x" + txHash
	url := fmt.Sprintf(ChainRPCJuno+"tx?hash=%s&prove=true", txHash)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error sending HTTP request: %s\n", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Handle the error
		fmt.Printf("Error reading response body: %s\n", err.Error())
		return nil, err
	}

	var data TxResultHttp
	if err := json.Unmarshal(body, &data); err != nil {
		// Handle the error
		fmt.Printf("Error unmarshalling JSON: %s\n", err.Error())
		return nil, err
	}

	return GetTxResult(&data), nil
}
func (Juno) GetNFT(classID, nftID string) (*NFT, error) { return nil, nil }
func (Juno) HasNFT(classID, nftID string) bool          { return false }
func (Juno) GetClass(classID string) (*Class, error)    { return nil, nil }
func (Juno) HasClass(classID string) bool               { return false }
