package chain

import (
	"golang.org/x/exp/slog"
	"io/ioutil"
	"net/http"
	"time"
)

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

func getRespWithRetry(url string) ([]byte, error) {
	var body []byte
	var resp *http.Response
	var err error
	maxRetries := 3

	for i := 1; i <= maxRetries; i++ {
		resp, err = http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			if i == maxRetries {
				return body, err
			}
			slog.Info("Http retry: ", i, " times")
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}
		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			if i == maxRetries {
				return body, err
			}
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}
		break
	}
	return body, err
}
