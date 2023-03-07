package chain

import (
	"context"
	"encoding/json"
	"fmt"
	nfttypes "github.com/UptickNetwork/uptick/x/collection/types"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
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

func (u Uptick) GetTx(txHash string) (*TxResult, error) {
	txHash = "0x" + txHash
	url := fmt.Sprintf(ChainRPCUptick+"tx?hash=%s&prove=true", txHash)

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

func (u Uptick) GetNFT(classID, nftID string) (*NFT, error) {
	req := &nfttypes.QueryNFTRequest{
		DenomId: classID,
		TokenId: nftID,
	}

	res, err := u.nftClient.NFT(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &NFT{
		ID:    res.NFT.Id,
		URI:   res.NFT.URI,
		Data:  res.NFT.Data,
		Owner: res.NFT.Owner,
	}, nil
}

func (u Uptick) HasNFT(classID, nftID string) bool {
	nft, _ := u.GetNFT(classID, nftID)
	if nft == nil {
		return false
	}
	return true
}

func (u Uptick) GetClass(classID string) (*Class, error) {
	req := nfttypes.QueryDenomRequest{
		DenomId: classID,
	}

	res, err := u.nftClient.Denom(context.Background(), &req)
	if err != nil {
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

func (u Uptick) HasClass(classID string) bool {
	nft, _ := u.GetClass(classID)
	if nft == nil {
		return false
	}
	return true
}
