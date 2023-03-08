package chain

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
