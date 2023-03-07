package verifier

import (
	"errors"
	"github.com/taramakage/gon-verifier/internal/chain"
)

type A3Params struct {
	ChainAbbreviation string
	TxHash            string
	ClassID           string // Wasm Contract Addr
	NFTID             string
	ChainID           string // Dest Chain ID
}

type A3Verifier struct {
	r *chain.Registry
}

func (v A3Verifier) Do(req Request, res chan<- *Response) {
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

	// FIXME: if A5 completed, owner is zero. if not, owner is sender.
	if nft.Owner != req.User.Address[params.ChainAbbreviation] || len(nft.Owner) != 0 {
		result.Reason = ReasonNFTOwnerNotMatch
		res <- result
		return
	}

	result.Point = PointMap[req.TaskNo]
	res <- result
}

func (v A3Verifier) BuildParams(rows [][]string) (any, error) {
	if len(rows) != 1 {
		return nil, errors.New("rows length not match")
	}

	param := rows[0]
	chainAbbr := chain.ChainIdAbbreviationStars
	if param[3] == chain.ChainIdValueJuno {
		chainAbbr = chain.ChainIdAbbreviationJuno
	}

	return A3Params{
		ChainAbbreviation: chainAbbr,
		TxHash:            param[0],
		ClassID:           param[1], // Wasm Contract Addr
		NFTID:             param[2],
		ChainID:           param[3],
	}, nil
}
