package chain

import (
	"context"
	"encoding/json"
	"fmt"
	nfttypes "github.com/UptickNetwork/uptick/x/collection/types"
	"github.com/taramakage/gon-verifier/internal/types"
	"google.golang.org/grpc"
)

type Uptick struct {
	conn      *grpc.ClientConn
	nftClient nfttypes.QueryClient
}

func NewUptick() *Uptick {
	conn, err := grpc.Dial(
		ChainGRPCUptick,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(),
	)
	if err != nil {
		panic(err)
	}

	return &Uptick{
		conn:      conn,
		nftClient: nfttypes.NewQueryClient(conn),
	}
}

func (u *Uptick) GetTx(txHash, txType string) (any, error) {
	txHash = "0x" + txHash
	url := fmt.Sprintf(ChainRPCUptick+"tx?hash=%s&prove=true", txHash)
	body, err := getRespWithRetry(url)
	if err != nil {
		return nil, err
	}

	var data types.TxResponse
	if err := json.Unmarshal(body, &data); err != nil {
		// Handle the error
		fmt.Printf("Error unmarshalling JSON: %s\n", err.Error())
		return nil, err
	}

	switch txType {
	case types.TxResultTypeIbcNft:
		return u.getTxResultIbcNft(&data)
	}

	return nil, fmt.Errorf("unknown tx type: %s", txType)
}

func (u *Uptick) getTxResultIbcNft(data *types.TxResponse) (any, error) {
	return data.IbcNftPkg()
}

func (u *Uptick) GetNFT(classID, nftID string) (*NFT, error) {
	req := &nfttypes.QueryNFTRequest{
		DenomId: classID,
		TokenId: nftID,
	}

	resi, err := withGrpcRetry(func() (interface{}, error) {
		return u.nftClient.NFT(context.Background(), req)
	})
	if err != nil {
		return nil, err
	}
	res, ok := resi.(*nfttypes.QueryNFTResponse)
	if !ok {
		return nil, err
	}

	return &NFT{
		ID:    res.NFT.Id,
		URI:   res.NFT.URI,
		Data:  res.NFT.Data,
		Owner: res.NFT.Owner,
	}, nil
}

func (u *Uptick) HasNFT(classID, nftID string) bool {
	nft, _ := u.GetNFT(classID, nftID)
	if nft == nil {
		return false
	}
	return true
}

func (u *Uptick) GetClass(classID string) (*Class, error) {
	req := &nfttypes.QueryDenomRequest{
		DenomId: classID,
	}

	resi, err := withGrpcRetry(func() (interface{}, error) {
		return u.nftClient.Denom(context.Background(), req)
	})
	if err != nil {
		return nil, err
	}
	res, ok := resi.(*nfttypes.QueryDenomResponse)
	if !ok {
		return nil, err
	}

	return &Class{
		ID:      res.Denom.Id,
		Name:    res.Denom.Name,
		Schema:  res.Denom.Schema,
		Creator: res.Denom.Creator,
		Uri:     res.Denom.Uri,
		UriHash: res.Denom.UriHash,
		Data:    res.Denom.Data,
	}, nil
}

func (u *Uptick) HasClass(classID string) bool {
	nft, _ := u.GetClass(classID)
	if nft == nil {
		return false
	}
	return true
}

func (u *Uptick) Close() {
	if u.conn != nil {
		u.conn.Close()
	}
}
