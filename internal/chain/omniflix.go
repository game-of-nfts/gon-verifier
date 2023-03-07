package chain

import (
	"context"
	"encoding/json"
	"fmt"
	nfttypes "github.com/OmniFlix/onft/types"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
)

type Omniflix struct {
	conn      *grpc.ClientConn
	nftClient nfttypes.QueryClient
}

func NewOmniflix() *Omniflix {
	conn, err := grpc.Dial(
		ChainGRPCIris,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(),
	)
	if err != nil {
		panic(err)
	}

	return &Omniflix{
		conn:      conn,
		nftClient: nfttypes.NewQueryClient(conn),
	}
}

func (o Omniflix) GetTx(txHash string) (*TxResult, error) {
	txHash = "0x" + txHash
	url := fmt.Sprintf(ChainRPCOmnilfix+"tx?hash=%s&prove=true", txHash)

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

func (o Omniflix) GetNFT(classID, nftID string) (*NFT, error) {
	req := &nfttypes.QueryONFTRequest{
		DenomId: classID,
		Id:      nftID,
	}

	res, err := o.nftClient.ONFT(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &NFT{
		ID:    res.ONFT.Id,
		URI:   res.ONFT.Metadata.PreviewURI, // NOTE
		Data:  res.ONFT.Data,
		Owner: res.ONFT.Owner,
	}, nil
}

func (o Omniflix) HasNFT(classID, nftID string) bool {
	nft, _ := o.GetNFT(classID, nftID)
	if nft == nil {
		return false
	}
	return true
}

func (o Omniflix) GetClass(classID string) (*Class, error) {
	req := nfttypes.QueryDenomRequest{
		DenomId: classID,
	}

	res, err := o.nftClient.Denom(context.Background(), &req)
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

func (o Omniflix) HasClass(classID string) bool {
	nft, _ := o.GetClass(classID)
	if nft == nil {
		return false
	}
	return true
}
