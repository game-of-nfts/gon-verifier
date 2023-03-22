package types

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

const (
	TxResultTypeRaw = "raw"
	TxResultTypeBasic      = "basic"
	TxResultTypeIssueDenom = "issue_denom"
	TxResultTypeMintNft    = "mint_nft"
	TxResultTypeIbcNft     = "ibc_nft"
)

type (
	TxResultBasic struct {
		Sender string
		TxCode int
	}

	TxResultIssueDenom struct {
		Sender  string
		Creator string
		DenomId string
		TxCode  int
	}

	TxResultMintNft struct {
		Sender    string
		DenomId   string
		TokenId   string
		Recipient string
		TxCode    int
	}

	TxResultIbcNft struct {
		Sender   string
		Receiver string
		DestPort string
		DestChan string
		ClassId  string
		TokenId  string
		TxCode   int
	}

	RaceResult struct {
		Sender   string
		Receiver string
		ClassId  string
		TokenId  string
		Height string
		TxCode   int
	}

	RaceData struct {
		Type      string `json:"type,omitempty"`
		Flow      string `json:"flow,omitempty"`
		LastOwner string `json:"last_owner,omitempty"`
		StartHeight string `json:"start_height,omitempty"`
		Name     struct {
			Value string `json:"value,omitempty"`
		} `json:"irismod:name,omitempty"`
		UriHash   struct {
				Value string `json:"value,omitempty"`
		} `json:"irismod:uri_hash,omitempty"`
	}

	IbcNftPacket struct {
		ClassId   string   `json:"classId"`
		ClassUri  string   `json:"classUri"`
		ClassData string   `json:"classData"`
		TokenIds  []string `json:"tokenIds"`
		TokenUris []string `json:"tokenUris"`
		TokenData []string `json:"tokenData"`
		Sender    string   `json:"sender"`
		Receiver  string   `json:"receiver"`
		Memo      string   `json:"memo"`
	}

	// TxResponse is the response of the tx query
	TxResponse struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  struct {
			Hash     string `json:"hash"`
			Height   string `json:"height"`
			Index    int    `json:"index"`
			TxResult struct {
				Code      int    `json:"code"`
				Data      string `json:"data"`
				Log       string `json:"log"`
				Info      string `json:"info"`
				GasWanted string `json:"gas_wanted"`
				GasUsed   string `json:"gas_used"`
				Events    []struct {
					Type       string `json:"type"`
					Attributes []struct {
						Key   string `json:"key"`
						Value string `json:"value"`
						Index bool   `json:"index"`
					} `json:"attributes"`
				} `json:"events"`
				Codespace string `json:"codespace"`
			} `json:"tx_result"`
			Tx    string `json:"tx"`
			Proof struct {
				RootHash string `json:"root_hash"`
				Data     string `json:"data"`
				Proof    struct {
					Total    string `json:"total"`
					Index    string `json:"index"`
					LeafHash string `json:"leaf_hash"`
					Aunts    []any  `json:"aunts"`
				} `json:"proof"`
			} `json:"proof"`
		} `json:"result"`
	}
)

// AttributeValueByKey  returns all the value of the given key in the given event.
func (tx *TxResponse) AttributeValueByKey(key string) string {
	kec := make([]byte, base64.StdEncoding.EncodedLen(len(key)))
	base64.StdEncoding.Encode(kec, []byte(key))
	for _, e := range tx.Result.TxResult.Events {
		for _, attr := range e.Attributes {
			// An event can have multiple same keys
			if attr.Key == string(kec) {
				value, _ := base64.StdEncoding.DecodeString(attr.Value)
				return string(value)
			}
		}
	}
	return ""
}

// EventAttributeValueByKey  returns all the value of the given key in the given event.
func (tx *TxResponse) EventAttributeValueByKey(event, key string) string {
	kec := make([]byte, base64.StdEncoding.EncodedLen(len(key)))
	base64.StdEncoding.Encode(kec, []byte(key))
	for _, e := range tx.Result.TxResult.Events {
		if e.Type == event {
			for _, attr := range e.Attributes {
				// An event can have multiple same keys
				if attr.Key == string(kec) {
					value, _ := base64.StdEncoding.DecodeString(attr.Value)
					return string(value)
				}
			}
		}
	}
	return ""
}

func (tx *TxResponse) DecodeAllEventAttributes() error {
	for _, e := range tx.Result.TxResult.Events {
		for i := range e.Attributes {
			key, _ := base64.StdEncoding.DecodeString(e.Attributes[i].Key)
			val, _ := base64.StdEncoding.DecodeString(e.Attributes[i].Value)
			e.Attributes[i].Key = string(key)
			e.Attributes[i].Value = string(val)
		}
	}
	return nil
}

func (tx *TxResponse) IbcNftPkg() (any, error) {
	ibcPkgRaw := tx.EventAttributeValueByKey(EventTypeIbcSendPacket, AttributeKeyIbcPackageData)
	var ibcPkg IbcNftPacket
	err := json.Unmarshal([]byte(ibcPkgRaw), &ibcPkg)
	if err != nil {
		return nil, err
	}

	return TxResultIbcNft{
		Sender:   ibcPkg.Sender,
		Receiver: ibcPkg.Receiver,
		DestPort: tx.EventAttributeValueByKey(EventTypeIbcSendPacket, AttributeKeyDestPort),
		DestChan: tx.EventAttributeValueByKey(EventTypeIbcSendPacket, AttributeKeyDestChan),
		ClassId:  ibcPkg.ClassId, // class-trace
		TokenId:  ibcPkg.TokenIds[0],
		TxCode:   tx.Result.TxResult.Code,
	}, nil
}

func (tx *TxResponse) GetFirstRace() (RaceResult, error) {
	ibcPkgRaw := tx.EventAttributeValueByKey(EventTypeIbcSendPacket, AttributeKeyIbcPackageData)
	var ibcPkg IbcNftPacket
	err := json.Unmarshal([]byte(ibcPkgRaw), &ibcPkg)
	if err != nil {
		return RaceResult{}, err
	}

	return RaceResult{
		Sender: ibcPkg.Sender,
		Receiver: ibcPkg.Receiver,
		ClassId: ibcPkg.ClassId,
		TokenId: ibcPkg.TokenIds[0],
		Height: tx.Result.Height,
		TxCode: tx.Result.TxResult.Code,
	}, err
}

func (tx *TxResponse) GetLastRace() (RaceResult, error) {
	return RaceResult{
		Sender:   tx.EventAttributeValueByKey(EventTypeNftTransfer, AttributeKeySender),
		Receiver: tx.EventAttributeValueByKey(EventTypeNftTransfer, AttributeKeyRecipient),
		ClassId:  tx.EventAttributeValueByKey(EventTypeNftTransfer, AttributeDenomId),
		TokenId:  tx.EventAttributeValueByKey(EventTypeNftTransfer, AttributeKeyTokenId),
		Height:   tx.Result.Height,
		TxCode:   tx.Result.TxResult.Code,
	}, nil
}

func (tx *TxResponse) GetIbcPkgRaceData() (RaceData, error) {
	ibcPkgRaw := tx.EventAttributeValueByKey(EventTypeIbcSendPacket, AttributeKeyIbcPackageData)
	var ibcPkg IbcNftPacket
	err := json.Unmarshal([]byte(ibcPkgRaw), &ibcPkg)
	if err != nil {
		return RaceData{}, err
	}
	if len(ibcPkg.TokenData) ==  0 {
		return RaceData{}, err
	}

	var raceData RaceData
	tokenData := ibcPkg.TokenData[0]
	decode, _ := base64.StdEncoding.DecodeString(tokenData)

	err = json.Unmarshal(decode, &raceData)
	if err != nil {
		return RaceData{}, err
	}
	return raceData, nil
}

func (txIbc *TxResultIbcNft) OriginalClass() string {
	ibcClassId := txIbc.ClassId
	elements := strings.Split(ibcClassId, "/")
	return elements[len(elements)-1]
}
