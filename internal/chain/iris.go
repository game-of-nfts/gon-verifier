package chain

import (
	"context"
	"encoding/json"
	"fmt"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"
	"github.com/taramakage/gon-verifier/internal/types"
	"google.golang.org/grpc"
)

type (
	Iris struct {
		conn      *grpc.ClientConn
		nftClient nfttypes.QueryClient
	}
)

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

// GetTx returns the transaction result
func (i *Iris) GetTx(txHash, txType string) (any, error) {
	txHash = "0x" + txHash
	url := fmt.Sprintf(ChainRPCIris+"tx?hash=%s&prove=true", txHash)
	body, err := getRespWithRetry(url)
	if err != nil {
		return nil, err
	}

	var data types.TxResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	switch txType {
	case types.TxResultTypeBasic:
		return i.getTxResultBasic(&data)
	case types.TxResultTypeIssueDenom:
		return i.getTxResultIssueDenom(&data)
	case types.TxResultTypeMintNft:
		return i.getTxResultMintNft(&data)
	case types.TxResultTypeIbcNft:
		return i.getTxResultIbcNft(&data)
	}

	return nil, fmt.Errorf("unknown tx type: %s", txType)
}

func (i *Iris) getTxResultBasic(data *types.TxResponse) (any, error) {
	return types.TxResultBasic{
		Sender: data.AttributeValueByKey(types.AttributeMsgSender),
		TxCode: data.Result.TxResult.Code,
	}, nil
}

func (i *Iris) getTxResultIssueDenom(data *types.TxResponse) (any, error) {
	return types.TxResultIssueDenom{
		Sender:  data.AttributeValueByKey(types.AttributeMsgSender),
		Creator: data.EventAttributeValueByKey(types.EventTypeIssueDenom, types.AttributeDenomCreator),
		DenomId: data.EventAttributeValueByKey(types.EventTypeIssueDenom, types.AttributeDenomId),
		TxCode:  data.Result.TxResult.Code,
	}, nil
}

func (i *Iris) getTxResultMintNft(data *types.TxResponse) (any, error) {
	return types.TxResultMintNft{
		Sender:    data.AttributeValueByKey(types.AttributeMsgSender),
		DenomId:   data.EventAttributeValueByKey(types.EventTypeNftMint, types.AttributeDenomId),
		TokenId:   data.EventAttributeValueByKey(types.EventTypeNftMint, types.AttributeKeyTokenId),
		Recipient: data.EventAttributeValueByKey(types.EventTypeNftMint, types.AttributeKeyRecipient),
		TxCode:    data.Result.TxResult.Code,
	}, nil
}

func (i *Iris) getTxResultIbcNft(data *types.TxResponse) (any, error) {
	return data.IbcNftPkg()
}

func (i *Iris) GetNFT(classID, nftID string) (*NFT, error) {
	req := &nfttypes.QueryNFTRequest{
		DenomId: classID,
		TokenId: nftID,
	}

	resi, err := withGrpcRetry(func() (interface{}, error) {
		return i.nftClient.NFT(context.Background(), req)
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

func (i *Iris) HasNFT(classID, nftID string) bool {
	nft, _ := i.GetNFT(classID, nftID)
	if nft == nil {
		return false
	}
	return true
}

func (i *Iris) GetClass(classID string) (*Class, error) {
	req := &nfttypes.QueryDenomRequest{
		DenomId: classID,
	}

	resi, err := withGrpcRetry(func() (interface{}, error) {
		return i.nftClient.Denom(context.Background(), req)
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

func (i *Iris) HasClass(classID string) bool {
	nft, _ := i.GetClass(classID)
	if nft == nil {
		return false
	}
	return true
}

func (i *Iris) Close() {
	if i.conn != nil {
		i.conn.Close()
	}
}
