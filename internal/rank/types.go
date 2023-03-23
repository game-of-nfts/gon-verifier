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

func BuildRaceInfo(r *RaceInfo, reason string) error {
	reason = strings.TrimSpace(reason)
	strs := strings.Split(reason, "/")
	if len(strs) != 4 {
		return errors.New("race format invalid")
	}
	start, err := strconv.Atoi(strs[1])
	if err != nil {
		return err
	}
	end, _ := strconv.Atoi(strs[2])
	if err != nil {
		return err
	}
	diff, _ := strconv.Atoi(strs[3])
	if err != nil {
		return err
	}

	r.start = start
	r.end = end
	r.diff = diff

	return nil
}
