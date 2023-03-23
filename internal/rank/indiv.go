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

type IndivRanker struct {
	TargetTaskNo   string
	TaskNo         string
	TaskPoint      int32
	IndivRaceInfos []IndivRaceInfo
	Entrance       string
}

func NewIndivRanker(entrance, targetTaskNo, taskNo string, taskPoint int32) IndivRanker {
	indivRaceInfos := make([]IndivRaceInfo, 0)
	return IndivRanker{
		TargetTaskNo:   targetTaskNo,
		TaskNo:         taskNo,
		TaskPoint:      taskPoint,
		IndivRaceInfos: indivRaceInfos,
		Entrance:       entrance,
	}
}

func (ir *IndivRanker) Do() error {
	err := ir.Collect()
	if err != nil {
		return err
	}
	ir.Sort()

	err = ir.GenerateRank()
	if err != nil {
		return err
	}

	err = ir.WriteTaskPoint()
	if err != nil {
		return err
	}
	return nil
}

// Collect will collect rank task results
func (ir *IndivRanker) Collect() error {
	files := make([]string, 0)
	err := filepath.Walk(ir.Entrance, func(path string, info os.FileInfo, err error) error {
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
		err = ir.loadRaceInfo(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ir *IndivRanker) loadRaceInfo(file string) error {
	taskpoint3, err := excelize.OpenFile(file)
	if err != nil {
		return err
	}
	defer taskpoint3.Close()

	rows, err := taskpoint3.GetRows("result")
	if err != nil {
		return errors.New("result sheet not found")
	}

	var indivRace *IndivRaceInfo = nil
	for _, row := range rows {
		if row[0] == ir.TargetTaskNo && strings.HasPrefix(row[3], "race") {
			raceInfo, err := BuildRaceInfo(row[3])
			if err != nil {
				return err
			}
			if raceInfo != nil {
				indivRace = &IndivRaceInfo{
					RaceInfo: *raceInfo,
					teamName: row[1],
					path:     file,
				}
			}
			break
		}
	}

	if indivRace != nil {
		ir.IndivRaceInfos = append(ir.IndivRaceInfos, *indivRace)
	}

	return nil
}

func (ir *IndivRanker) Sort() {
	sort.SliceStable(ir.IndivRaceInfos, func(i, j int) bool {
		ri1 := ir.IndivRaceInfos[i].RaceInfo
		ri2 := ir.IndivRaceInfos[j].RaceInfo
		if ri1.diff == ri2.diff {
			return ri1.start < ri2.start
		}
		return ri1.diff < ri2.diff
	})
}

func (ir *IndivRanker) GenerateRank() error {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "result"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	f.SetCellValue(sheetName, "A1", "Rank")
	f.SetCellValue(sheetName, "B1", "TeamName")
	f.SetCellValue(sheetName, "C1", "DiffHeight")
	f.SetCellValue(sheetName, "D1", "StartHeight")
	f.SetCellValue(sheetName, "E1", "EndHeight")

	for i, indivRaceInfo := range ir.IndivRaceInfos {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), indivRaceInfo.teamName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), indivRaceInfo.diff)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), indivRaceInfo.start)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), indivRaceInfo.end)
	}

	f.SetActiveSheet(index)

	fileName := filepath.Join(ir.Entrance, fmt.Sprintf("rank%s.xlsx", ir.TaskNo))
	if err := f.SaveAs(fileName); err != nil {
		return err
	}
	return nil
}

func (ir *IndivRanker) WriteTaskPoint() error {
	for i, indivRaceInfo := range ir.IndivRaceInfos {
		// top 10
		if i == 10 {
			break
		}
		err := ir.clearLegacyTaskPoint(indivRaceInfo)
		if err != nil {
			return err
		}
		err = ir.writeTaskPoint(indivRaceInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ir *IndivRanker) clearLegacyTaskPoint(indivRaceInfo IndivRaceInfo) error {
	file, err := excelize.OpenFile(indivRaceInfo.path)
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
		if rows[i][0] == ir.TaskNo {
			file.RemoveRow(sheetName, i+1)
		}
	}

	err = file.Save()
	if err != nil {
		return err
	}

	return nil
}

func (ir *IndivRanker) writeTaskPoint(indivRaceInfo IndivRaceInfo) error {
	file, err := excelize.OpenFile(indivRaceInfo.path)
	if err != nil {
		return err
	}
	defer file.Close()

	sheetName := "result"
	rows, err := file.GetRows(sheetName)
	length := len(rows)

	file.SetCellValue(sheetName, fmt.Sprintf("A%d", length+1), ir.TaskNo)
	file.SetCellValue(sheetName, fmt.Sprintf("B%d", length+1), indivRaceInfo.teamName)
	file.SetCellValue(sheetName, fmt.Sprintf("C%d", length+1), ir.TaskPoint)

	err = file.Save()
	if err != nil {
		return err
	}

	return nil
}
