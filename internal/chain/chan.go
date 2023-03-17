package chain

import (
	"errors"
	"strings"
)

var PortChanPairStrMap = map[string]string{
	"is-1": "nft-transfer/channel-22 <> wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-207",
	"is-2": "nft-transfer/channel-23 <> wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-208",
	"ij-1": "nft-transfer/channel-24 <> wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-89",
	"ij-2": "nft-transfer/channel-25 <> wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-90",
	"iu-1": "nft-transfer/channel-17 <> nft-transfer/channel-3",
	"iu-2": "nft-transfer/channel-19 <> nft-transfer/channel-4",
	"io-1": "nft-transfer/channel-0 <> nft-transfer/channel-24",
	"io-2": "nft-transfer/channel-1 <> nft-transfer/channel-25",
	"sj-1": "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-211 <> wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-93",
	"sj-2": "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-213 <> wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-94",
	"sj-3": "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-230 <> wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-120",
	"sj-4": "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-234 <> wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-122",
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
	id   string
	src  PortChan
	dest PortChan
}

func NewPortChanPair(src, dest, chanPairId string) (*PortChanPair, error) {
	key := src + dest + "-" + chanPairId
	if _, ok := PortChanPairStrMap[key]; ok {
		value := PortChanPairStrMap[key]
		parts := strings.Split(value, " <> ")
		spc := strings.Split(parts[0], "/")
		dpc := strings.Split(parts[1], "/")
		return &PortChanPair{
			id: chanPairId,
			src: PortChan{
				ChainAbbr: src,
				Port:      spc[0],
				Channel:   spc[1],
			},
			dest: PortChan{
				ChainAbbr: dest,
				Port:      dpc[0],
				Channel:   dpc[1],
			},
		}, nil
	}

	keyr := dest + src + "-" + chanPairId
	if _, ok := PortChanPairStrMap[keyr]; ok {
		value := PortChanPairStrMap[keyr]
		parts := strings.Split(value, " <> ")
		dpc := strings.Split(parts[0], "/")
		spc := strings.Split(parts[1], "/")
		return &PortChanPair{
			id: chanPairId,
			src: PortChan{
				ChainAbbr: src,
				Port:      spc[0],
				Channel:   spc[1],
			},
			dest: PortChan{
				ChainAbbr: dest,
				Port:      dpc[0],
				Channel:   dpc[1],
			},
		}, nil
	}

	return nil, errors.New("port channel pair not found")
}

func (p *PortChanPair) IsEquivalent(p2 *PortChanPair) bool {
	if p.id == p2.id && p.src.ChainAbbr == p2.dest.ChainAbbr && p.dest.ChainAbbr == p2.src.ChainAbbr {
		return true
	}
	return false
}

func (p *PortChanPair) GetSrcPortChan() PortChan {
	return p.src
}

func (p *PortChanPair) GetDestPortChan() PortChan {
	return p.dest
}
