package verifier

import (
	"encoding/json"
	"errors"
	"github.com/taramakage/gon-verifier/internal/chain"
	"github.com/taramakage/gon-verifier/internal/types"
)

type A1Params struct {
	ChainAbbreviation string
	TxHash            string
	ClassId           string
}

type A1Verifier struct {
	r *chain.Registry
}

type A1ClassData struct {
	GithubUsername string `json:"github_username"`
	DiscordHandle  string `json:"discord_handle,omitempty"`
	TeamName       string `json:"team_name,omitempty"`
	Community      string `json:"community,omitempty"`
}

func (v A1Verifier) Do(req Request, res chan<- *Response) {
	result := &Response{
		TaskNo:   req.TaskNo,
		TeamName: req.User.TeamName,
	}

	// params validation
	params, ok := req.Params.(A1Params)
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

	c := v.r.GetChain(params.ChainAbbreviation)
	txi, err := c.GetTx(params.TxHash, types.TxResultTypeIssueDenom)
	if err != nil {
		result.Reason = ReasonTxResultUnachievable
		res <- result
		return
	}
	tx, ok := txi.(types.TxResultIssueDenom)
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

	if req.User.Address[params.ChainAbbreviation] != tx.Sender {
		result.Reason = ReasonTxMsgSenderNotMatch
		res <- result
		return
	}

	// query class on chain
	class, err := c.GetClass(params.ClassId)
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

	if len(class.Uri) == 0 {
		result.Reason = ReasonClassUrIEmpty
		res <- result
		return
	}

	var classData A1ClassData
	err = json.Unmarshal([]byte(class.Data), &classData)
	if err != nil {
		result.Reason = ReasonClassDataInvalid
		res <- result
		return
	}

	result.Point = PointMap[req.TaskNo]
	res <- result
}

func (v A1Verifier) BuildParams(rows [][]string) (any, error) {
	if len(rows) < 1 {
		return nil, errors.New("format is incorrect")
	}
	rowFirst := rows[0]
	return A1Params{
		ChainAbbreviation: chain.ChainIdAbbreviationIris,
		TxHash:            rowFirst[0],
		ClassId:           rowFirst[1],
	}, nil
}
