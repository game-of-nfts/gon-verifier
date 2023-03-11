package chain

import (
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
	flow    string
	chains  []rune
	chanIds []rune
	maxHop  int
	curr    int
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

	return &Flow{
		flow:    flowString,
		chains:  chains,
		chanIds: chanIds,
		maxHop:  len(chanIds),
		curr:    0,
	}, nil
}

// Next returns the next transfer info for a flow
func (f *Flow) Next() (*PortChanPair, bool) {
	if f.curr == f.maxHop {
		return &PortChanPair{}, false
	}

	id := f.chanIds[f.curr]
	c1 := f.chains[f.curr]
	c2 := f.chains[f.curr+1]
	pcp, err := NewPortChanPair(string(c1), string(c2), string(id))
	if err != nil {
		return &PortChanPair{}, false
	}
	f.curr++

	return pcp, true
}

func (f *Flow) Done() bool {
	return f.curr == f.maxHop
}

func (f *Flow) CalcIbcClassHash(original string, flowId string) string {
	// TODO
	return ""
}
