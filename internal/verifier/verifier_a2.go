package verifier

import (
	"fmt"
	"github.com/taramakage/gon-verifier/internal/chain"
	"github.com/taramakage/gon-verifier/internal/types"
	"strings"
)

type A2Params struct {
	ChainAbbreviation string
	TxHashes          []string
	ClassIds          []string
	TokenIds          []string
	ParamErrorMsg     string
}

type A2Verifier struct {
	r *chain.Registry
}

func (v A2Verifier) Do(req Request, res chan<- *Response) {
	result := &Response{
		TaskNo:   req.TaskNo,
		TeamName: req.User.TeamName,
	}

	// params validation
	params, ok := req.Params.(A2Params)
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
	if len(params.ChainAbbreviation) == 0 {
		result.Reason = ReasonParamsChainIdEmpty
		res <- result
		return
	}

	// a2 validation
	c := v.r.GetChain(params.ChainAbbreviation)
	for i := range params.TxHashes {
		txi, err := c.GetTx(params.TxHashes[i], types.TxResultTypeMintNft)
		if err != nil {
			result.Reason = ReasonTxResultUnachievable
			res <- result
			return
		}
		tx, ok := txi.(types.TxResultMintNft)
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

		// class owner must be the same as register address on iris
		class, err := c.GetClass(params.ClassIds[i])
		if err != nil {
			result.Reason = ReasonClassNotFound
			res <- result
			return
		}
		if req.User.Address[params.ChainAbbreviation] != class.Creator {
			result.Reason = ReasonClassCreatorNotMatch
			res <- result
			return
		}

		if req.User.Address[params.ChainAbbreviation] != tx.Sender {
			result.Reason = ReasonTxMsgSenderNotMatch
			res <- result
			return
		}

		if req.User.Address[params.ChainAbbreviation] != tx.Recipient {
			result.Reason = ReasonNftOwnerNotMatch
			res <- result
			return
		}

		// query nft on chain
		nft, err := c.GetNFT(params.ClassIds[i], params.TokenIds[i])
		if err != nil {
			result.Reason = ReasonNftNotFound
			res <- result
			return
		}

		if len(nft.URI) == 0 {
			result.Reason = ReasonNftUriEmpty
			res <- result
			return
		}

		if len(nft.Data) == 0 {
			result.Reason = ReasonNftDataEmpty
			res <- result
			return
		}
	}

	result.Point = PointMap[req.TaskNo]
	res <- result
}

func (v A2Verifier) BuildParams(rows [][]string) (any, error) {
	if len(rows) < 2 {
		return A2Params{
			ParamErrorMsg: fmt.Sprintf("parmas of task wanted at least 2 row(s) , but got %d row(s)", len(rows)),
		}, nil
	} else if strings.HasPrefix(strings.TrimSpace(rows[0][0]), "tx") {
		return A2Params{
			ParamErrorMsg: "row 2 should be replaced with evidence rather than left there",
		}, nil
	}

	params := A2Params{
		ChainAbbreviation: chain.ChainIdAbbreviationIris,
		TxHashes:          make([]string, 0),
		ClassIds:          make([]string, 0),
		TokenIds:          make([]string, 0),
	}

	// NOTE: only the first two rows are read
	for i := range rows {
		if i == 2 {
			break
		}

		params.TxHashes = append(params.TxHashes, rows[i][0])
		params.ClassIds = append(params.ClassIds, rows[i][1])
		params.TokenIds = append(params.TokenIds, rows[i][2])
	}

	return params.Trim(), nil
}

func (p A2Params) Trim() A2Params {
	res := p
	for i := range res.TxHashes {
		res.TxHashes[i] = strings.TrimSpace(res.TxHashes[i])
		res.ClassIds[i] = strings.TrimSpace(res.ClassIds[i])
		res.TokenIds[i] = strings.TrimSpace(res.TokenIds[i])
	}
	return res
}
