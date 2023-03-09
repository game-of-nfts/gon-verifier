package chain

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
			time.Sleep(time.Duration(i*2) * time.Second)
			continue
		}
		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			if i == maxRetries {
				fmt.Printf("Http write body times: %d\n", maxRetries)
				return body, err
			}
			fmt.Printf("Http write body times: %d\n", i)
			time.Sleep(time.Duration(i*2) * time.Second)
			continue
		}
		break
	}
	return body, err
}

func withGrpcRetry(fn func() (interface{}, error)) (interface{}, error) {
	var (
		retryCount   = 0
		maxRetries   = 3
		retryBackoff = time.Second
	)

	for {
		res, err := fn()
		if err != nil {
			if status.Code(err) == codes.Unavailable || status.Code(err) == codes.DeadlineExceeded {
				if retryCount >= maxRetries {
					return nil, fmt.Errorf("max grpc retries reached")
				}
				retryCount++
				time.Sleep(time.Duration(retryCount) * retryBackoff)
				continue
			} else {
				return nil, err
			}
		} else {
			return res, nil
		}
	}
}
