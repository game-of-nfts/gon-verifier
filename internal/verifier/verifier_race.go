package verifier

import (
	"fmt"
	"github.com/taramakage/gon-verifier/internal/chain"
	"github.com/taramakage/gon-verifier/internal/types"
	"strconv"
	"strings"
)

type RaceParam struct {
	firstTransfer string
	lastTransfer  string
	ParamErrorMsg string
}

// RaceVerifier validates whether a participant has completed task B1,B2,B5,B6,B7.
type RaceVerifier struct {
	r *chain.Registry
	f  *chain.Flow
	originalClassId string
	designatedOwner string
	startBlockHeight string
	endBlockHeight string
}

func NewRaceVerifier(r *chain.Registry, originalClassId string, designatedOwner string, startBlockHeight, endBlockHeight string) *RaceVerifier {
	// flow init is delayed to verify
	return &RaceVerifier{
		r:   r,
		originalClassId: originalClassId,
		designatedOwner: designatedOwner,
		startBlockHeight: startBlockHeight,
		endBlockHeight: endBlockHeight,
	}
}

func (v RaceVerifier) Do(req Request, res chan<- *Response) {
	result := &Response{
		TaskNo:   req.TaskNo,
		TeamName: req.User.TeamName,
	}

	params, ok := req.Params.(RaceParam)
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

	iris := v.r.GetChain(chain.ChainIdAbbreviationIris)
	txi1, err := iris.GetTx(params.firstTransfer, types.TxResultTypeRaw)
	if err != nil {
		result.Reason = ReasonTxResultUnachievable
		res <- result
		return
	}
	tx1, ok := txi1.(types.TxResponse)
	if !ok {
		result.Reason = ReasonTxResultUnexpected
		res <- result
		return
	}

	// build flow according to flow-id
	race, _ := tx1.GetIbcPkgRaceData()
	f, _ := chain.NewFlow(chain.FlowStrMap[race.Flow])
	v.f = f

	txi2, err := iris.GetTx(params.lastTransfer, types.TxResultTypeRaw)
	if err != nil {
		result.Reason = ReasonTxResultUnachievable
		res <- result
		return
	}
	tx2, ok := txi2.(types.TxResponse)
	if !ok {
		result.Reason = ReasonTxResultUnexpected
		res <- result
		return
	}

	first, _ := tx1.GetFirstRace()
	last, _ := tx2.GetLastRace()

	hash, _ := v.f.GetFinalIbcHash(v.originalClassId)
	ibcClass := "ibc/" + hash.String()
	if ibcClass != last.ClassId {
		result.Reason = ReasonRaceUnexpectedFlowPath
		res <- result
		return
	}

	if first.Sender != last.Sender {
		result.Reason = ReasonRaceFirstLastSenderNotMatch
		res <- result
		return
	}

	if first.Sender !=  req.User.Address[chain.ChainIdAbbreviationIris] {
		result.Reason = ReasonTxMsgSenderNotMatch
		res <- result
		return
	}

	nft, err := iris.GetNFT(last.ClassId, last.TokenId)
	if err != nil {
		result.Reason = ReasonNftNotFound
		res <- result
		return
	}

	if nft.Owner != v.designatedOwner {
		result.Reason = ReasonNftOwnerNotMatch
		res <- result
		return
	}

	if race.StartHeight < v.startBlockHeight {
		result.Reason = ReasonRaceStartTooEarly
		res <- result
		return
	}

	result.Point = PointMap[req.TaskNo]
	if last.Height <= v.endBlockHeight {
		result.Reason = v.BuildRaceResult(first.Height, last.Height)
	}

	res <- result
}

func (v RaceVerifier) BuildRaceResult(first, last string) string {
	l, _ := strconv.Atoi(last)
	f, _ := strconv.Atoi(first)
	diff := l - f
	return fmt.Sprintf("race/%s/%s/%s", first, last, strconv.Itoa(diff))
}

func (v RaceVerifier) BuildParams(rows [][]string) (any, error) {
	errMsg := restrictParamLen(rows, 2)
	if len(errMsg) != 0 {
		return RaceParam{
			ParamErrorMsg: errMsg,
		}, nil
	}

	params := RaceParam{
		firstTransfer: rows[0][0],
		lastTransfer:  rows[1][0],
	}
	return params.Trim(), nil
}

func (p RaceParam) Trim() RaceParam {
	res := p
	res.firstTransfer = strings.TrimSpace(res.firstTransfer)
	res.lastTransfer = strings.TrimSpace(res.lastTransfer)
	return res
}
