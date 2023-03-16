package verifier

import (
	"github.com/taramakage/gon-verifier/internal/chain"
	"github.com/taramakage/gon-verifier/internal/types"
	"strings"
)

type FlowParams struct {
	TxHashes        []string
	IbcClassId      string // ibc/hash on iris
	OriginalClassId string // original class id
	TokenId         string // token-id
	ParamErrorMsg   string
}

type FlowVerifier struct {
	r   *chain.Registry
	f   *chain.Flow
	ngb bool
}

func NewFlowVerifier(r *chain.Registry, flowId string, ngb bool) *FlowVerifier {
	flowStr, ok := chain.FlowStrMap[flowId]
	if !ok {
		return nil
	}
	f, err := chain.NewFlow(flowStr)
	if err != nil {
		return nil
	}
	return &FlowVerifier{
		r:   r,
		f:   f,
		ngb: ngb,
	}
}

func (v FlowVerifier) Do(req Request, res chan<- *Response) {
	result := &Response{
		TaskNo:   req.TaskNo,
		TeamName: req.User.TeamName,
	}

	params, ok := req.Params.(FlowParams)
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

	if !v.ngb {
		if ok, reason := v.ValidateByTxHash(&params, &req); !ok {
			result.Reason = reason
			res <- result
			return
		}
	}

	if ok, reason := v.ValidateByIbcClass(&params, &req); !ok {
		result.Reason = reason
		res <- result
		return
	}

	result.Point = PointMap[req.TaskNo]
	res <- result
}

// ValidateByIbcClass check the owner of nft under ibc class on last destination
func (v FlowVerifier) ValidateByIbcClass(param *FlowParams, req *Request) (bool, string) {
	// check nft existence
	iris := v.r.GetChain(chain.ChainIdAbbreviationIris)
	nft, err := iris.GetNFT(param.IbcClassId, param.TokenId)
	if err != nil {
		return false, ReasonNftNotFound
	}
	// check owner of nft
	if req.User.Address[chain.ChainIdAbbreviationIris] != nft.Owner {
		return false, ReasonNftOwnerNotMatch
	}
	// ibc class trace match the flow
	hash, _ := v.f.GetFinalIbcHash(param.OriginalClassId)
	ibc := "ibc/" + hash.String()
	if ibc != param.IbcClassId {
		return false, ReasonIbcClassNotMatch
	}
	return true, ""
}

// ValidateByTxHash validate each tx hash according the flow
func (v FlowVerifier) ValidateByTxHash(param *FlowParams, req *Request) (bool, string) {
	for i, txHash := range param.TxHashes {
		// get tx result
		srcChain := v.r.GetChain(v.f.GetSrcChainAbbr(i))
		txi, err := srcChain.GetTx(txHash, types.TxResultTypeIbcNft)
		if err != nil {
			return false, ReasonTxResultUnachievable
		}
		tx, ok := txi.(types.TxResultIbcNft)
		if !ok {
			return false, ReasonTxResultUnexpected
		}
		if tx.TxCode != 0 {
			return false, ReasonTxResultUnsuccessful
		}

		pcp := v.f.GetPortChanPairByIdx(i)
		dpc := pcp.GetDestPortChan()

		if tx.DestPort != dpc.Port {
			return false, ReasonIbcDestPortNotMatch
		}
		if tx.DestChan != dpc.Channel {
			return false, ReasonIbcDestChanNotMatch
		}
		if tx.Sender != req.User.Address[v.f.GetSrcChainAbbr(i)] {
			return false, ReasonTxMsgSenderNotMatch
		}
		if tx.Receiver != req.User.Address[v.f.GetDestChainAbbr(i)] {
			return false, ReasonNftRecipientNotMatch
		}
		if tx.TokenId != param.TokenId {
			return false, ReasonNftTokenIdNotMatch
		}
	}

	return true, ""
}

func (v FlowVerifier) BuildParams(rows [][]string) (any, error) {
	if v.ngb {
		return v.buildParamsNgb(rows)
	}
	return v.buildParams(rows)
}

// buildParamsNgb build params from never-go-back transfer evidence
// - txHashes: nil
// - ibcClassId: provided by rows
// - tokenId: provided by rows
// - originalClassId:
func (v FlowVerifier) buildParamsNgb(rows [][]string) (any, error) {
	errMsg := restrictParamLen(rows, 1)
	if len(errMsg) != 0 {
		return FlowParams{
			ParamErrorMsg: errMsg,
		}, nil
	}

	params := FlowParams{
		TxHashes:      nil,
		IbcClassId:    rows[0][0],
		TokenId:       rows[0][1],
		ParamErrorMsg: "",
	}
	return params.Trim().AddOriginalClassId(&v), nil
}

// buildParams build params from non never-go-back transfer evidence
// - txHashes: provided by rows
// - ibcClassId: calculated by flow-id and the first txHash
// - tokenId: calculated until the first txHash is used
func (v FlowVerifier) buildParams(rows [][]string) (any, error) {
	maxHop := v.f.GetFlowHops()
	errMsg := restrictParamLen(rows, maxHop)
	if len(errMsg) != 0 {
		return FlowParams{
			ParamErrorMsg: errMsg,
		}, nil
	}

	params := FlowParams{
		TxHashes: make([]string, maxHop),
	}
	for i := range rows {
		params.TxHashes[i] = rows[i][0]
	}

	return params.Trim().AddThreeKindId(&v), nil
}

func (p FlowParams) Trim() FlowParams {
	res := p
	res.TokenId = strings.TrimSpace(res.TokenId)
	res.IbcClassId = strings.TrimSpace(res.IbcClassId)
	for i := range res.TxHashes {
		res.TxHashes[i] = strings.TrimSpace(res.TxHashes[i])
	}
	return res
}

func (p FlowParams) AddOriginalClassId(v *FlowVerifier) FlowParams {
	irisi := v.r.GetChain(chain.ChainIdAbbreviationIris)
	iris, ok := irisi.(*chain.Iris)
	if !ok {
		p.ParamErrorMsg = ReasonParamsFormatIncorrect
		return p
	}
	originalClassId, err := iris.GetOriginalClassId(p.IbcClassId)
	if err != nil {
		p.ParamErrorMsg = ReasonIbcOriginalClassIdNotMatch
		return p
	}
	p.OriginalClassId = originalClassId
	return p
}

func (p FlowParams) AddThreeKindId(v *FlowVerifier) FlowParams {
	// get tx result
	srcChain := v.r.GetChain(chain.ChainIdAbbreviationIris)
	txi, err := srcChain.GetTx(p.TxHashes[0], types.TxResultTypeIbcNft)
	if err != nil {
		p.ParamErrorMsg = ReasonTxResultUnachievable
		return p
	}
	tx, ok := txi.(types.TxResultIbcNft)
	if !ok {
		p.ParamErrorMsg = ReasonTxResultUnexpected
		return p
	}
	if tx.TxCode != 0 {
		p.ParamErrorMsg = ReasonTxResultUnsuccessful
		return p
	}
	p.TokenId = tx.TokenId
	p.OriginalClassId = tx.OriginalClass()
	// NOTE: ibc class id is not provided by user, so we need to calculate it
	hash, _ := v.f.GetFinalIbcHash(p.OriginalClassId)
	p.IbcClassId = "ibc/" + hash.String()
	return p
}
