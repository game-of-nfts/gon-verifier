package chain

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/taramakage/gon-verifier/internal/types"
	wasmtype "github.com/taramakage/gon-verifier/internal/types/wasm"
	"google.golang.org/grpc"
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

func (j Juno) GetTx(txHash, txType string) (any, error) {
	txHash = "0x" + txHash
	url := fmt.Sprintf(ChainRPCJuno+"tx?hash=%s&prove=true", txHash)
	body, err := getRespWithRetry(url)
	if err != nil {
		return nil, err
	}

	var data types.TxResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	switch txType {
	case types.TxResultTypeIbcNft:
		return j.getTxResultIbcNft(&data)
	}

	return nil, fmt.Errorf("unknown tx type: %s", txType)
}

func (j Juno) getTxResultIbcNft(data *types.TxResponse) (any, error) {
	return data.IbcNftPkg()
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
