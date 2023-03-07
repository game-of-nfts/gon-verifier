package chain

import "encoding/base64"

type WasmQueryNFT struct {
	NftInfo NftInfo `json:"nft_info"`
}
type NftInfo struct {
	TokenId string `json:"token_id"`
}

type WasmRespNFT struct {
}

type WasmQueryClass struct {
	NumTokens NumTokens `json:"num_tokens"`
}
type NumTokens struct {
}

type WasmRespClass struct {
	Count int `json:"count"`
}

func GetTxResult(data *TxResultHttp) *TxResult {
	txResult := TxResult{}
	for _, event := range data.Result.TxResult.Events {
		if event.Type == "message" {
			for _, attr := range event.Attributes {
				key, _ := base64.StdEncoding.DecodeString(attr.Key)
				if string(key) == "sender" {
					value, _ := base64.StdEncoding.DecodeString(attr.Value)
					txResult.Sender = string(value)
					return &txResult
				}
			}
		}
	}
	return nil
}
