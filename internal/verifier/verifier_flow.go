package verifier

//
//import "github.com/taramakage/gon-verifier/internal/chain"
//
//type FlowParams struct {
//	IsNGB         bool
//	TxHash        []string
//	ChainId       []string
//	IbcClassId    string
//	TokenId       string
//	ParamErrorMsg string
//}
//
//type FlowVerifier struct {
//	r *chain.Registry
//}
//
//func (v FlowVerifier) Do(req Request, res chan<- *Response) {
//	result := &Response{
//		TaskNo:   req.TaskNo,
//		TeamName: req.User.TeamName,
//	}
//
//}
//
//func (v FlowVerifier) ValidateByIbcClass(param *FlowParams) {
//
//}
//
//func (v FlowVerifier) ValidateByHash(param *FlowParams) {
//
//}
//
//func (p FlowParams) Validation(res chan<- *Response) {
//
//}
//
//func (p FlowParams) Trim() FlowParams {
//	return FlowParams{
//		IsNGB:         p.IsNGB,
//		TxHash:        p.TxHash,
//		ChainId:       p.ChainId,
//		IbcClassId:    p.IbcClassId,
//		TokenId:       p.TokenId,
//		ParamErrorMsg: p.ParamErrorMsg,
//	}
//}
