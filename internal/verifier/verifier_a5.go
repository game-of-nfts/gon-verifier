package verifier

import (
	"github.com/taramakage/gon-verifier/internal/chain"
	"github.com/taramakage/gon-verifier/internal/types"
	"strings"
)

type A5Params struct {
	ChainAbbreviation string
	TxHash            string
	ClassId           string // Wasm Contract Addr
	TokenId           string
	ChainId           string // Dest Chain ID
	ParamErrorMsg     string
}

type A5Verifier struct {
	r *chain.Registry
}

func (v A5Verifier) Do(req Request, res chan<- *Response) {
	result := &Response{
		TaskNo:   req.TaskNo,
		TeamName: req.User.TeamName,
	}

	params, ok := req.Params.(A5Params)
	if !ok {
		result.Reason = ReasonParamsFormatIncorrect
		res <- result
		return
	}
	if len(params.ParamErrorMsg) != 0 {
		result.Reason = params.ParamErrorMsg
		res <- result
		return
	}
	if len(params.TxHash) == 0 {
		result.Reason = ReasonParamsChainIdEmpty
		res <- result
		return
	}
	if len(params.ChainAbbreviation) == 0 {
		result.Reason = ReasonParamsChainIdError
		res <- result
		return
	}

	srcChain := v.r.GetChain(params.ChainAbbreviation)
	txi, err := srcChain.GetTx(params.TxHash, types.TxResultTypeIbcNft)
	if err != nil {
		result.Reason = ReasonTxResultUnachievable
		res <- result
		return
	}
	tx, ok := txi.(types.TxResultIbcNft)
	if !ok {
		result.Reason = ReasonTxResultUnexpected
		res <- result
		return
	}
	if tx.TxCode != 0 {
		result.Reason = ReasonTxResultUnsuccessful
		res <- result
		return
	}

	// query cw-721 addr on chain
	if !srcChain.HasClass(params.ClassId) {
		result.Reason = ReasonClassNotFound
		res <- result
		return
	}

	if req.User.Address[params.ChainAbbreviation] != tx.Sender {
		result.Reason = ReasonTxMsgSenderNotMatch
		res <- result
		return
	}
	if req.User.Address[chain.ChainIdAbbreviationIris] != tx.Receiver {
		result.Reason = ReasonNftRecipientNotMatch
		res <- result
		return
	}

	iris := v.r.GetChain(chain.ChainIdAbbreviationIris)
	originalClassId := tx.OriginalClass()
	if !iris.HasNFT(originalClassId, params.TokenId) {
		result.Reason = ReasonNftNotFound
		res <- result
		return
	}

	result.Point = PointMap[req.TaskNo]
	res <- result
}

func (v A5Verifier) BuildParams(rows [][]string) (any, error) {
	errMsg := restrictParamLen(rows, 1)
	if len(errMsg) != 0 {
		return A5Params{
			ParamErrorMsg: errMsg,
		}, nil
	}

	param := rows[0]
	return A5Params{
		ChainAbbreviation: "",
		TxHash:            param[0],
		ClassId:           param[1], // Wasm Contract Addr
		TokenId:           param[2],
		ChainId:           param[3],
	}.Trim(), nil
}

func (p A5Params) Trim() A5Params {
	res := p
	res.TxHash = strings.TrimSpace(res.TxHash)
	res.ClassId = strings.TrimSpace(res.ClassId)
	res.ClassId = strings.TrimPrefix(res.ClassId, "wasm.")
	res.TokenId = strings.TrimSpace(res.TokenId)
	res.ChainId = strings.TrimSpace(res.ChainId)

	if res.ChainId == chain.ChainIdValueJuno {
		res.ChainAbbreviation = chain.ChainIdAbbreviationJuno
	}
	if res.ChainId == chain.ChainIdValueStars {
		res.ChainAbbreviation = chain.ChainIdAbbreviationStars
	}
	return res
}
