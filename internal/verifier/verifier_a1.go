package verifier

import (
	"errors"

	"github.com/taramakage/gon-verifier/internal/chain"
)

type A1Params struct {
	ChainAbbreviation string
	TxHash            string
	ClassID           string
}

type A1Verifier struct {
	r *chain.Registry
}

func (v A1Verifier) Do(req Request, res chan<- *Response) {
	result := &Response{
		TaskNo:   req.TaskNo,
		TeamName: req.User.TeamName,
	}

	params, ok := req.Params.(A1Params)
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

	chain := v.r.GetChain(params.ChainAbbreviation)
	tx, err := chain.GetTx(params.TxHash)
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

	class, err := chain.GetClass(params.ClassID)
	if err != nil {
		result.Reason = err.Error()
		res <- result
		return
	}

	if class.Creator != req.User.Address[params.ChainAbbreviation] {
		result.Reason = ReasonClassCreatorNotMatch
		res <- result
		return
	}

	if len(class.Uri) == 0 {
		result.Reason = ReasonClassURIEmpty
		res <- result
		return
	}

	// TODO: validate class data content.
	if len(class.Data) == 0 {
		result.Reason = ReasonClassDataEmpty
		res <- result
		return
	}

	result.Point = PointMap[req.TaskNo]
	res <- result
}

func (v A1Verifier) BuildParams(rows [][]string) (any, error) {
	// FIXME: only the first row is read
	if len(rows) != 1 {
		return nil, errors.New("非法的格式，只能提交一行数据")
	}
	rowFirst := rows[0]
	return A1Params{
		ChainAbbreviation: chain.ChainIdAbbreviationIris,
		TxHash:            rowFirst[0],
		ClassID:           rowFirst[1],
	}, nil
}
