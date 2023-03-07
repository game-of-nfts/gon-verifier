package chain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Stargaze struct {
}

func (Stargaze) GetTx(txHash string) (*TxResult, error) {
	txHash = "0x" + txHash
	url := fmt.Sprintf(ChainRPCStars+"tx?hash=%s&prove=true", txHash)

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
func (Stargaze) GetNFT(classID, nftID string) (*NFT, error) {
	return nil, nil
}
func (Stargaze) HasNFT(classID, nftID string) bool       { return false }
func (Stargaze) GetClass(classID string) (*Class, error) { return nil, nil }
func (Stargaze) HasClass(classID string) bool            { return false }
