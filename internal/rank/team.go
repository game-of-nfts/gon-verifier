package rank

import (
	"errors"
	"fmt"
	"github.com/taramakage/gon-verifier/internal/scorecard"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type TeamRanker struct {
	TargetTaskNos []string
	TaskNo        string
	TaskPoint     int32
	TeamRaceInfos []TeamRaceInfo
	Entrance      string
}

func NewTeamRanker(entrance string, targetTaskNos []string, taskNo string, taskPoint int32) TeamRanker {
	teamRaceInfos := make([]TeamRaceInfo, 0)
	return TeamRanker{
		TargetTaskNos: targetTaskNos,
		TaskNo:        taskNo,
		TaskPoint:     taskPoint,
		TeamRaceInfos: teamRaceInfos,
		Entrance:      entrance,
	}
}

func (tr *TeamRanker) Do() error {
	err := tr.Collect()
	if err != nil {
		return err
	}
	tr.Sort()

	err = tr.GenerateRank()
	if err != nil {
		return err
	}

	err = tr.WriteTaskPoint()
	if err != nil {
		return err
	}
	return nil
}

func (tr *TeamRanker) Collect() error {
	files := make([]string, 0)
	err := filepath.Walk(tr.Entrance, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() || info.Name() != scorecard.DefaultStageThreeTaskPoint {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return err
	}

	for _, file := range files {
		err := tr.loadRaceInfo(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (tr *TeamRanker) loadRaceInfo(file string) error {
	taskpoint3, err := excelize.OpenFile(file)
	if err != nil {
		return err
	}
	defer taskpoint3.Close()

	rows, err := taskpoint3.GetRows("result")
	if err != nil {
		return errors.New("result sheet not found")
	}

	teamRace := TeamRaceInfo{
		raceInfos: make([]RaceInfo, 0),
	}

	for _, row := range rows {
		for _, targetTaskNo := range tr.TargetTaskNos {
			if row[0] == targetTaskNo && strings.HasPrefix(row[3], "race") {
				raceInfo, err := BuildRaceInfo(row[3])
				if err != nil {
					return err
				}
				if raceInfo != nil {
					teamRace.raceInfos = append(teamRace.raceInfos, *raceInfo)
				}
			}
		}
	}

	if len(teamRace.raceInfos) != 0 {
		teamRace.teamName = rows[1][1]
		teamRace.path = file
		for _, raceInfo := range teamRace.raceInfos {
			teamRace.diffSum += raceInfo.diff
			teamRace.startSum += raceInfo.start
		}

		tr.TeamRaceInfos = append(tr.TeamRaceInfos, teamRace)
	}

	return nil
}

func (tr *TeamRanker) Sort() {
	sort.SliceStable(tr.TeamRaceInfos, func(i, j int) bool {
		diff1 := tr.TeamRaceInfos[i].diffSum
		diff2 := tr.TeamRaceInfos[j].diffSum
		if diff1 == diff2 {
			return tr.TeamRaceInfos[i].startSum < tr.TeamRaceInfos[i].startSum
		}
		return diff1 < diff2
	})
}

func (tr *TeamRanker) GenerateRank() error {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "result"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	f.SetCellValue(sheetName, "A1", "Rank")
	f.SetCellValue(sheetName, "B1", "TeamName")
	f.SetCellValue(sheetName, "C1", "SumOfDiffHeight")
	f.SetCellValue(sheetName, "D1", "SumOfStartHeight")

	for i, teamRaceInfo := range tr.TeamRaceInfos {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), teamRaceInfo.teamName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), teamRaceInfo.diffSum)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), teamRaceInfo.startSum)
	}

	f.SetActiveSheet(index)

	fileName := filepath.Join(tr.Entrance, fmt.Sprintf("rank%s.xlsx", tr.TaskNo))
	if err := f.SaveAs(fileName); err != nil {
		return err
	}
	return nil
}

func (tr *TeamRanker) WriteTaskPoint() error {
	for i, teamRaceInfo := range tr.TeamRaceInfos {
		// top 10
		if i == 10 {
			break
		}
		err := tr.clearLegacyTaskPoint(teamRaceInfo)
		if err != nil {
			return err
		}
		err = tr.writeTaskPoint(teamRaceInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tr *TeamRanker) clearLegacyTaskPoint(teamRaceInfo TeamRaceInfo) error {
	file, err := excelize.OpenFile(teamRaceInfo.path)
	if err != nil {
		return err
	}
	defer file.Close()

	sheetName := "result"
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return err
	}

	// remove in reverse order to avoid changing the index of the remaining rows
	for i := len(rows) - 1; i >= 0; i-- {
		if rows[i][0] == tr.TaskNo {
			file.RemoveRow(sheetName, i+1)
		}
	}

	err = file.Save()
	if err != nil {
		return err
	}

	return nil
}

func (tr *TeamRanker) writeTaskPoint(teamRaceInfo TeamRaceInfo) error {
	file, err := excelize.OpenFile(teamRaceInfo.path)
	if err != nil {
		return err
	}
	defer file.Close()

	sheetName := "result"
	rows, err := file.GetRows(sheetName)
	length := len(rows)

	file.SetCellValue(sheetName, fmt.Sprintf("A%d", length+1), tr.TaskNo)
	file.SetCellValue(sheetName, fmt.Sprintf("B%d", length+1), teamRaceInfo.teamName)
	file.SetCellValue(sheetName, fmt.Sprintf("C%d", length+1), tr.TaskPoint)

	err = file.Save()
	if err != nil {
		return err
	}

	return nil
}