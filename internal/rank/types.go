package rank

import (
	"errors"
	"strconv"
	"strings"
)

type RaceInfo struct {
	start int
	end   int
	diff  int
}

type IndivRaceInfo struct {
	RaceInfo
	teamName string
	path     string
}

type TeamRaceInfo struct {
	raceInfos []RaceInfo
	diffSum   int
	startSum  int
	teamName  string
	path      string
}

func NewTeamRaceInfo() *TeamRaceInfo {
	raceInfos := make([]RaceInfo, 0)
	return &TeamRaceInfo{
		raceInfos: raceInfos,
	}
}

func BuildRaceInfo(reason string) (*RaceInfo, error) {
	reason = strings.TrimSpace(reason)
	strs := strings.Split(reason, "/")
	if len(strs) != 4 {
		return nil, errors.New("race format invalid")
	}
	start, err := strconv.Atoi(strs[1])
	if err != nil {
		return nil, err
	}
	end, _ := strconv.Atoi(strs[2])
	if err != nil {
		return nil, err
	}
	diff, _ := strconv.Atoi(strs[3])
	if err != nil {
		return nil, err
	}

	return &RaceInfo{
		start: start,
		end:   end,
		diff:  diff,
	}, nil
}
