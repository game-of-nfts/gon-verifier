package chain

import (
	"crypto/sha256"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	"unicode"
)

const (
	FlowA01 = "i --(1)--> s --(1)--> j --(1)--> i"
	FlowA02 = "i --(1)--> u --(1)--> o --(1)--> i"
	FlowA03 = "i --(1)--> s --(1)--> j --(1)--> u --(1)--> i"
	FlowA04 = "i --(1)--> s --(1)--> o --(1)--> j --(1)--> i"
	FlowA05 = "i --(1)--> s --(1)--> j --(1)--> u --(1)--> o --(1)--> s --(1)--> i"
	FlowA06 = "i --(1)--> o --(1)--> s --(1)--> u --(1)--> o --(1)--> j --(1)--> i"
	FlowB01 = "i --(1)--> s --(1)--> u --(1)--> s --(2)--> i"
	FlowB02 = "i --(1)--> u --(1)--> o --(1)--> u --(2)--> i"
	FlowB03 = "i --(1)--> j --(1)--> u --(1)--> j --(2)--> i"
	FlowB04 = "i --(1)--> j --(1)--> s --(1)--> j --(2)--> i"
	FlowC01 = "i --(1)--> s --(1)--> j --(1)--> s --(1)--> i"
	FlowC02 = "i --(1)--> o --(1)--> u --(1)--> o --(1)--> i"
	FlowC03 = "i --(1)--> s --(1)--> j --(1)--> u --(1)--> j --(1)--> s --(1)--> i"
	FlowC04 = "i --(1)--> u --(1)--> s --(1)--> o --(1)--> s --(1)--> u --(1)--> i"
)

var FlowStrMap = map[string]string{
	"a01": FlowA01,
	"a02": FlowA02,
	"a03": FlowA03,
	"a04": FlowA04,
	"a05": FlowA05,
	"a06": FlowA06,
	"b01": FlowB01,
	"b02": FlowB02,
	"b03": FlowB03,
	"b04": FlowB04,
	"c01": FlowC01,
	"c02": FlowC02,
	"c03": FlowC03,
	"c04": FlowC04,
}

type Flow struct {
	flow      string
	transfers [][3]rune // ["src-abbr", "dest-abbr", "pc-id"]
	maxHop    int
	curr      int
}

func NewFlow(flowString string) (*Flow, error) {
	chanIds := make([]rune, 0)
	chains := make([]rune, 0)
	for _, c := range flowString {
		if unicode.IsDigit(c) {
			chanIds = append(chanIds, c)
		}
		if unicode.IsLetter(c) {
			chains = append(chains, c)
		}
	}

	transfers := make([][3]rune, 0)
	for i, c := range chanIds {
		transfer := [3]rune{chains[i], chains[i+1], c}
		transfers = append(transfers, transfer)
	}

	return &Flow{
		flow:      flowString,
		transfers: transfers,
		maxHop:    len(chanIds),
		curr:      0,
	}, nil
}

func (f *Flow) GetFlowHops() int {
	return f.maxHop
}

// Next returns the next transfer port/chan info in a flow
func (f *Flow) Next() (*PortChanPair, bool) {
	if f.curr == f.maxHop {
		return nil, false
	}

	transfer := f.transfers[f.curr]

	pcp, err := NewPortChanPair(string(transfer[0]), string(transfer[1]), string(transfer[2]))
	if err != nil {
		return nil, false
	}
	f.curr++

	return pcp, true
}

func (f *Flow) Done() bool {
	return f.curr == f.maxHop
}

func (f *Flow) GetSrcChainAbbr(idx int) string {
	if idx >= f.maxHop {
		return ""
	}
	return string(f.transfers[idx][0])
}

func (f *Flow) GetDestChainAbbr(idx int) string {
	if idx >= f.maxHop {
		return ""
	}
	return string(f.transfers[idx][1])
}

// GetFinalIbcHash calculates hash of trace/classId
func (f *Flow) GetFinalIbcHash(classId string) (tmbytes.HexBytes, error) {
	trace, err := f.buildFinalClassTrace()
	if err != nil {
		return nil, err
	}
	if len(classId) != 0 {
		classId = trace + classId
	}
	hash := sha256.Sum256([]byte(classId))
	return hash[:], nil
}

func (f *Flow) GetOriginalHash(classId string) (tmbytes.HexBytes, error) {
	hash := sha256.Sum256([]byte(classId))
	return hash[:], nil
}

// GetFinalClassTrace returns the final ibc class trace in the flow
func (f *Flow) GetFinalClassTrace() (string, error) {
	return f.buildFinalClassTrace()
}

func (f *Flow) GetPortChanPair(transfer [3]rune) *PortChanPair {
	pcp, err := NewPortChanPair(string(transfer[0]), string(transfer[1]), string(transfer[2]))
	if err != nil {
		return nil
	}
	return pcp
}

func (f *Flow) GetPortChanPairByIdx(idx int) *PortChanPair {
	transfer := [3]rune{f.transfers[idx][0], f.transfers[idx][1], f.transfers[idx][2]}
	return f.GetPortChanPair(transfer)
}

// buildFinalClassTrace returns the final ibc class trace in the flow
// (a -> b)                 p1/c1/class
// a -> (b -> c)            p2/c2/p1/c1/class
// a -> b -> (c -> b)       p1/c1/class
// a -> b -> c -> (b -> d)  p3/c3/p1/c1/class
func (f *Flow) buildFinalClassTrace() (string, error) {
	transferTrim := make([][3]rune, 0)
	for i, transfer := range f.transfers {
		if i == 0 {
			transferTrim = append(transferTrim, f.transfers[i])
			continue
		}
		k := len(transferTrim)
		if transferTrim[k-1][2] == transfer[2] && transferTrim[k-1][0] == transfer[1] && transferTrim[k-1][1] == transfer[0] {
			transferTrim = transferTrim[:k-1]
		} else {
			transferTrim = append(transferTrim, f.transfers[i])
		}
	}

	trace := ""
	for _, t := range transferTrim {
		pcp := f.GetPortChanPair(t)
		trace = string(pcp.dest.Port) + "/" + string(pcp.dest.Channel) + "/" + trace
	}
	return trace, nil
}
