package chain

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/taramakage/gon-verifier/internal/types"
	wasmtype "github.com/taramakage/gon-verifier/internal/types/wasm"
	"google.golang.org/grpc"
)

type Stargaze struct {
	conn       *grpc.ClientConn
	wasmClient wasmtype.QueryClient
}

func NewStargaze() *Stargaze {
	conn, err := grpc.Dial(
		ChainGRPCStars,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(),
	)
	if err != nil {
		panic(err)
	}

	return &Stargaze{
		conn:       conn, // NOTE: Close this connection when the program exits
		wasmClient: wasmtype.NewQueryClient(conn),
	}
}

func (s *Stargaze) GetTx(txHash, txType string) (any, error) {
	txHash = "0x" + txHash
	url := fmt.Sprintf(ChainRPCStars+"tx?hash=%s&prove=true", txHash)
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
		return s.getTxResultIbcNft(&data)
	}

	return nil, fmt.Errorf("unknown tx type: %s", txType)
}

func (s *Stargaze) getTxResultIbcNft(data *types.TxResponse) (any, error) {
	return data.IbcNftPkg()
}

func (s *Stargaze) GetNFT(classID, nftID string) (*NFT, error) {
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
	resi, err := withGrpcRetry(func() (interface{}, error) {
		return s.wasmClient.SmartContractState(context.Background(), req)
	})
	if err != nil {
		return nil, err
	}
	_, ok := resi.(*wasmtype.QuerySmartContractStateResponse)
	if !ok {
		return nil, err
	}

	return &NFT{
		ID: nftID,
	}, nil
}

func (s *Stargaze) HasNFT(classID, nftID string) bool {
	_, err := s.GetNFT(classID, nftID)
	if err != nil {
		return false
	}
	return true
}

func (s *Stargaze) GetClass(classID string) (*Class, error) {
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
	resi, err := withGrpcRetry(func() (interface{}, error) {
		return s.wasmClient.SmartContractState(context.Background(), req)
	})
	if err != nil {
		return nil, err
	}
	res, ok := resi.(*wasmtype.QuerySmartContractStateResponse)
	if !ok {
		return nil, err
	}

	wr := WasmRespClass{}
	err = json.Unmarshal(res.Data, &wr)
	if err != nil {
		return nil, nil
	}

	return &Class{}, nil
}

func (s *Stargaze) HasClass(classID string) bool {
	_, err := s.GetClass(classID)
	if err != nil {
		return false
	}
	return true
}

func (s *Stargaze) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}
