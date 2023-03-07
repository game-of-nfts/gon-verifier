package verifier

import (
	"errors"
	"github.com/taramakage/gon-verifier/internal/chain"
)

type A2Params struct {
	ChainAbbreviation string
	TxHash            []string
	ClassID           []string
	NFTID             []string
}

type A2Verifier struct {
	r *chain.Registry
}

func (v A2Verifier) Do(req Request, res chan<- *Response) {
	result := &Response{
		TaskNo:   req.TaskNo,
		TeamName: req.User.TeamName,
	}

	params, ok := req.Params.(A2Params)
	if !ok {
		result.Reason = ReasonIllegalParams
		res <- result
		return
	}

	if len(params.ChainAbbreviation) == 0 {
		result.Reason = ReasonChainIDEmpty
		res <- result
		return
	}

	chain := v.r.GetChain(params.ChainAbbreviation)
	for i := range params.TxHash {
		tx, err := chain.GetTx(params.TxHash[i])
		if err != nil {
			result.Reason = err.Error()
			res <- result
			return
		}

		if req.User.Address[params.ChainAbbreviation] != tx.Sender {
			result.Reason = ReasonSenderNotMatch
			res <- result
			return
		}

		nft, err := chain.GetNFT(params.ClassID[i], params.NFTID[i])
		if err != nil {
			result.Reason = err.Error()
			res <- result
			return
		}

		// FIXME: add escrow address logic
		if nft.Owner != req.User.Address[params.ChainAbbreviation] {
			result.Reason = ReasonNFTOwnerNotMatch
			res <- result
			return
		}

		if len(nft.URI) == 0 {
			result.Reason = ReasonNFTURIEmpty
			res <- result
			return
		}

		if len(nft.Data) == 0 {
			result.Reason = ReasonNFTDataEmpty
			res <- result
			return
		}
	}

	result.Point = PointMap[req.TaskNo]
	res <- result
}

func (v A2Verifier) BuildParams(rows [][]string) (any, error) {
	if len(rows) < 2 {
		return nil, errors.New("参数行数不足")
	}

	params := A2Params{
		ChainAbbreviation: chain.ChainIdAbbreviationIris,
		TxHash:            make([]string, 0),
		ClassID:           make([]string, 0),
		NFTID:             make([]string, 0),
	}

	// NOTE: only the first two rows are read
	for i := range rows {
		if i == 0 {
			continue
		}

		if i == 3 {
			break
		}

		params.TxHash = append(params.TxHash, rows[i][0])
		params.ClassID = append(params.ClassID, rows[i][1])
		params.NFTID = append(params.NFTID, rows[i][2])
	}

	return params, nil
}
