package verifier

import (
	"errors"
	"github.com/taramakage/gon-verifier/internal/chain"
	"github.com/taramakage/gon-verifier/internal/types"
)

type A5Params struct {
	ChainAbbreviation string
	TxHash            string
	ClassId           string // Wasm Contract Addr
	TokenId           string
	ChainId           string // Dest Chain ID
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
	if len(rows) != 1 {
		return nil, errors.New("task evidence format is incorrect")
	}

	param := rows[0]
	chainAbbr := ""
	if param[3] == chain.ChainIdValueJuno {
		chainAbbr = chain.ChainIdAbbreviationJuno
	}
	if param[3] == chain.ChainIdValueStars {
		chainAbbr = chain.ChainIdAbbreviationStars
	}

	return A5Params{
		ChainAbbreviation: chainAbbr,
		TxHash:            param[0],
		ClassId:           param[1], // Wasm Contract Addr
		TokenId:           param[2],
		ChainId:           param[3],
	}, nil
}
