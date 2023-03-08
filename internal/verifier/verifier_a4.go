package verifier

import (
	"errors"
	"github.com/taramakage/gon-verifier/internal/chain"
	"github.com/taramakage/gon-verifier/internal/types"
)

type A4Params struct {
	ChainAbbreviation string
	TxHash            string
	ClassId           string // IBC Class
	TokenId           string
	ChainId           string // Dest Chain Id
}

type A4Verifier struct {
	r *chain.Registry
}

func (v A4Verifier) Do(req Request, res chan<- *Response) {
	result := &Response{
		TaskNo:   req.TaskNo,
		TeamName: req.User.TeamName,
	}

	params, ok := req.Params.(A4Params)
	if !ok {
		result.Reason = ReasonParamsFormatIncorrect
		res <- result
		return
	}
	if len(params.ChainAbbreviation) == 0 {
		result.Reason = ReasonParamsChainIdEmpty
		res <- result
		return
	}

	iris := v.r.GetChain(chain.ChainIdAbbreviationIris)
	txi, err := iris.GetTx(params.TxHash, types.TxResultTypeIbcNft)
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

	// query ibc class on chain
	destChain := v.r.GetChain(params.ChainAbbreviation)
	if ok := destChain.HasClass(params.ClassId); !ok {
		result.Reason = ReasonClassNotFound
		res <- result
		return
	}

	if req.User.Address[chain.ChainIdAbbreviationIris] != tx.Sender {
		result.Reason = ReasonTxMsgSenderNotMatch
		res <- result
		return
	}
	if req.User.Address[params.ChainAbbreviation] != tx.Receiver {
		result.Reason = ReasonNftRecipientNotMatch
		res <- result
		return
	}

	if tx.TokenId != params.TokenId {
		result.Reason = ReasonNftTokenIdNotMatch
		res <- result
		return
	}

	result.Point = PointMap[req.TaskNo]
	res <- result
}

func (v A4Verifier) BuildParams(rows [][]string) (any, error) {
	if len(rows) < 1 {
		return nil, errors.New("format is incorrect")
	}

	param := rows[0]
	chainAbbr := ""
	if param[3] == chain.ChainIdValueOmniflix {
		chainAbbr = chain.ChainIdAbbreviationOmniflix
	}
	if param[3] == chain.ChainIdValueUptick {
		chainAbbr = chain.ChainIdAbbreviationUptick
	}

	return A4Params{
		ChainAbbreviation: chainAbbr,
		TxHash:            param[0],
		ClassId:           param[1],
		TokenId:           param[2],
		ChainId:           param[3],
	}, nil
}
