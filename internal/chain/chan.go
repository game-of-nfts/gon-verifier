package chain

import (
	"errors"
	"fmt"
	"strings"
)

var PortChanPairStrMap = map[string]string{
	"is-1": "nft-transfer/channel-22 <> wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-207",
	"is-2": "nft-transfer/channel-23 <> wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-208",
	"ij-1": "nft-transfer/channel-24 <> wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-89",
	"ij-2": "nft-transfer/channel-25 <> wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-90",
	"iu-1": "nft-transfer/channel-17 <> nft-transfer/channel-3",
	"iu-2": "nft-transfer/channel-19 <> nft-transfer/channel-4",
	"io-1": "nft-transfer/channel-0 <> nft-transfer/channel-24t",
	"io-2": "nft-transfer/channel-1 <> nft-transfer/channel-25",
	"sj-1": "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-211 <> wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-93",
	"sj-2": "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-213 <> wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-94",
	"su-1": "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-203 <> nft-transfer/channel-6",
	"su-2": "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-206 <> nft-transfer/channel-12",
	"so-1": "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-209 <> nft-transfer/channel-44",
	"so-2": "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-210 <> nft-transfer/channel-45",
	"ju-1": "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-86 <> nft-transfer/channel-7",
	"ju-2": "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-88 <> nft-transfer/channel-13",
	"jo-1": "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-91 <> nft-transfer/channel-46",
	"jo-2": "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-92 <> nft-transfer/channel-47",
	"uo-1": "nft-transfer/channel-5 <> nft-transfer/channel-41",
	"uo-2": "nft-transfer/channel-9 <> nft-transfer/channel-42",
}

type PortChan struct {
	ChainAbbr string
	Port      string
	Channel   string
}

type PortChanPair struct {
	id  string
	cp1 PortChan
	cp2 PortChan
}

func NewPortChanPair(chainAbbr1, chainAbbr2, chanPairId string) (*PortChanPair, error) {
	key := chainAbbr1 + chainAbbr2 + "-" + chanPairId
	value := ""
	if _, ok := PortChanPairStrMap[key]; ok {
		value = PortChanPairStrMap[key]
	} else {
		tmp := chainAbbr1
		chainAbbr1 = chainAbbr2
		chainAbbr2 = tmp

		key = chainAbbr1 + chainAbbr2 + "-" + chanPairId
		if _, ok := PortChanPairStrMap[key]; ok {
			value = PortChanPairStrMap[key]
		}
	}

	if len(value) == 0 {
		return nil, errors.New(fmt.Sprintf("unknown port chain pair: %s", key))
	}

	parts := strings.Split(value, " <> ")
	pc1 := strings.Split(parts[0], "/")
	pc2 := strings.Split(parts[1], "/")

	return &PortChanPair{
		id: chanPairId,
		cp1: PortChan{
			chainAbbr1,
			pc1[0],
			pc1[1],
		},
		cp2: PortChan{
			chainAbbr2,
			pc2[0],
			pc2[1],
		},
	}, nil
}
