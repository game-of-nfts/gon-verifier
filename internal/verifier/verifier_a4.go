package verifier

import (
	"errors"
	"github.com/taramakage/gon-verifier/internal/chain"
)

type A4Params struct {
	ChainAbbreviation string
	TxHash            string
	ClassID           string // Wasm Contract Addr
	NFTID             string
	ChainID           string // Dest Chain ID
}

type A4Verifier struct {
	r *chain.Registry
}

func (v A4Verifier) Do(req Request, res chan<- *Response) {
	result := &Response{
		TaskNo:   req.TaskNo,
		TeamName: req.User.TeamName,
	}

	params, ok := req.Params.(A3Params)
	if !ok {
		result.Reason = ReasonIllegalParams
		res <- result
		return
	}

	if len(params.TxHash) == 0 {
		result.Reason = ReasonTxHashEmpty
		res <- result
		return
	}

	// NOTE: check tx on iris
	iris := v.r.GetChain(chain.ChainIdAbbreviationIris)
	tx, err := iris.GetTx(params.TxHash)
	if err != nil {
		result.Reason = err.Error()
		res <- result
		return
	}

	if req.User.Address[chain.ChainIdAbbreviationIris] != tx.Sender {
		result.Reason = ReasonSenderNotMatch
		res <- result
		return
	}

	destChain := v.r.GetChain(params.ChainAbbreviation)
	nft, err := destChain.GetNFT(params.ClassID, params.NFTID)
	if err != nil {
		result.Reason = err.Error()
		res <- result
		return
	}

	// FIXME: if A6 completed, owner is empty. if not, owner is user.
	if nft.Owner != req.User.Address[params.ChainAbbreviation] || len(nft.Owner) != 0 {
		result.Reason = ReasonNFTOwnerNotMatch
		res <- result
		return
	}

	result.Point = PointMap[req.TaskNo]
	res <- result
}

func (v A4Verifier) BuildParams(rows [][]string) (any, error) {
	if len(rows) != 1 {
		return nil, errors.New("rows length not match")
	}

	param := rows[1]
	chainAbbr := chain.ChainIdAbbreviationUptick
	if param[3] == chain.ChainIdValueOmniflix {
		chainAbbr = chain.ChainIdAbbreviationOmniflix
	}

	return A4Params{
		ChainAbbreviation: chainAbbr,
		TxHash:            param[0],
		ClassID:           param[1],
		NFTID:             param[2],
		ChainID:           param[3],
	}, nil
}
