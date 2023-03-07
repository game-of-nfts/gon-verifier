package chain

import (
	"context"
	"encoding/json"
	"fmt"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
)

type Iris struct {
	conn      *grpc.ClientConn
	nftClient nfttypes.QueryClient
}

func NewIris() *Iris {
	conn, err := grpc.Dial(
		ChainGRPCIris,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(),
	)
	if err != nil {
		panic(err)
	}

	return &Iris{
		conn:      conn, // NOTE: Close this connection when the program exits
		nftClient: nfttypes.NewQueryClient(conn),
	}
}

func (i Iris) GetTx(txHash string) (*TxResult, error) {
	txHash = "0x" + txHash
	url := fmt.Sprintf(ChainRPCIris+"tx?hash=%s&prove=true", txHash)

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

func (i Iris) GetNFT(classID, nftID string) (*NFT, error) {
	req := &nfttypes.QueryNFTRequest{
		DenomId: classID,
		TokenId: nftID,
	}

	res, err := i.nftClient.NFT(context.Background(), req)
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

func (i Iris) HasNFT(classID, nftID string) bool {
	nft, _ := i.GetNFT(classID, nftID)
	if nft == nil {
		return false
	}
	return true
}

func (i Iris) GetClass(classID string) (*Class, error) {
	req := nfttypes.QueryDenomRequest{
		DenomId: classID,
	}

	res, err := i.nftClient.Denom(context.Background(), &req)
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

func (i Iris) HasClass(classID string) bool {
	nft, _ := i.GetClass(classID)
	if nft == nil {
		return false
	}
	return true
}
