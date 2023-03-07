package chain

import (
	"context"
	"encoding/json"
	"fmt"
	wasmtype "github.com/taramakage/gon-verifier/x/wasm/types"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
)

type Juno struct {
	conn       *grpc.ClientConn
	wasmClient wasmtype.QueryClient
}

func NewJuno() *Juno {
	conn, err := grpc.Dial(
		ChainGRPCJuno,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(),
	)
	if err != nil {
		panic(err)
	}

	return &Juno{
		conn:       conn, // NOTE: Close this connection when the program exits
		wasmClient: wasmtype.NewQueryClient(conn),
	}
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
func (j Juno) GetNFT(classID, nftID string) (*NFT, error) {
	wq := WasmQueryNFT{
		NftInfo: NftInfo{nftID},
	}
	// convert wq to json string
	bz, err := json.Marshal(wq)
	if err != nil {
		return nil, err
	}

	req := &wasmtype.QuerySmartContractStateRequest{
		Address:   classID,
		QueryData: bz,
	}
	_, err = j.wasmClient.SmartContractState(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &NFT{
		ID: nftID,
	}, nil
}

func (j Juno) HasNFT(classID, nftID string) bool {
	_, err := j.GetNFT(classID, nftID)
	if err != nil {
		return false
	}
	return true
}
func (j Juno) GetClass(classID string) (*Class, error) {
	wq := WasmQueryClass{}
	// convert wq to json string
	bz, err := json.Marshal(wq)
	if err != nil {
		return nil, err
	}

	req := &wasmtype.QuerySmartContractStateRequest{
		Address:   classID,
		QueryData: bz,
	}
	res, err := j.wasmClient.SmartContractState(context.Background(), req)
	fmt.Println(res)
	if err != nil {
		return nil, err
	}

	wr := WasmRespClass{}
	err = json.Unmarshal(res.Data, &wr)
	if err != nil {
		return nil, nil
	}

	return &Class{}, nil
}

func (j Juno) HasClass(classID string) bool {
	_, err := j.GetClass(classID)
	if err != nil {
		return false
	}
	return true
}
