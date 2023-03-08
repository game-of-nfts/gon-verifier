package scorecard

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	DefaultScoreCardFile  = "scorecard.xlsx"
	DefaultTaskPointFile  = "taskpoint.xlsx"
	DefaultTaskPointSheet = "result"
	DefaultScoreCardSheet = "result"
)

// ScoreCard reads task results and output the scorecard
type ScoreCard struct {
	entranceDir string // path: entrance
}

type ScoreCardEntry struct {
	taskCompleted string
	totalPoint    int
	teamName      string
}

// NewScoreCard creates a new ScoreCard
func NewScoreCard(entranceDir string) *ScoreCard {
	return &ScoreCard{entranceDir: entranceDir}
}

func (sc *ScoreCard) Generate() error {
	f := excelize.NewFile()
	index, err := f.NewSheet(DefaultScoreCardSheet)
	if err != nil {
		return err
	}
	f.SetCellValue(DefaultScoreCardSheet, "A1", "rank")
	f.SetCellValue(DefaultScoreCardSheet, "B1", "team_name")
	f.SetCellValue(DefaultScoreCardSheet, "C1", "task_completed")
	f.SetCellValue(DefaultScoreCardSheet, "D1", "final_score")
	f.SetCellValue(DefaultScoreCardSheet, "E1", "update_time")

	var taskPointFiles []string
	err = filepath.Walk(sc.entranceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == DefaultTaskPointFile {
			taskPointFiles = append(taskPointFiles, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	var entries []*ScoreCardEntry
	for _, taskPointFile := range taskPointFiles {
		entry, err := sc.HandleTaskPoint(taskPointFile)
		if err != nil {
			fmt.Println(taskPointFiles, " ", err)
			continue
		}
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].totalPoint == entries[j].totalPoint {
			return entries[i].teamName < entries[j].teamName
		}

		return entries[i].totalPoint > entries[j].totalPoint
	})

	for i := 0; i < len(entries); i++ {
		f.SetCellValue(DefaultScoreCardSheet, "A"+strconv.Itoa(i+2), i+1)
		f.SetCellValue(DefaultScoreCardSheet, "B"+strconv.Itoa(i+2), entries[i].teamName)
		f.SetCellValue(DefaultScoreCardSheet, "C"+strconv.Itoa(i+2), entries[i].taskCompleted)
		f.SetCellValue(DefaultScoreCardSheet, "D"+strconv.Itoa(i+2), entries[i].totalPoint)
	}

	f.SetActiveSheet(index)
	scorecardFile := filepath.Join(sc.entranceDir, DefaultScoreCardFile)
	if err := f.SaveAs(scorecardFile); err != nil {
		return err
	}

	return nil
}

// HandleTaskPoint receives a taskpoint.xlsx and returns its scorecard entry
func (sc *ScoreCard) HandleTaskPoint(taskPointFile string) (*ScoreCardEntry, error) {
	taskpoint, err := excelize.OpenFile(taskPointFile)
	if err != nil {
		return nil, err
	}

	rows, err := taskpoint.GetRows(DefaultTaskPointSheet)
	if len(rows) <= 1 {
		return nil, err
	}

	var taskResults TaskResults
	for _, row := range rows[1:] {
		point, _ := strconv.Atoi(row[2])
		taskResults = append(taskResults, TaskResult{
			TaskNo: row[0],
			Point:  point,
		})
	}

	return &ScoreCardEntry{
		taskCompleted: sc.concatenateTaskNo(taskResults),
		totalPoint:    sc.calculateTotalPoint(taskResults),
		teamName:      rows[1][1],
	}, nil
}

func (sc *ScoreCard) concatenateTaskNo(taskResults TaskResults) string {
	sort.Sort(taskResults)
	taskNos := make([]string, 0)
	for i := 0; i < len(taskResults); i++ {
		if taskResults[i].Point != 0 {
			taskNos = append(taskNos, taskResults[i].TaskNo)
		}
	}

	return strings.Join(taskNos, ",")
}

func (sc *ScoreCard) calculateTotalPoint(taskResults TaskResults) int {
	var totalPoint int
	for _, taskResult := range taskResults {
		totalPoint += taskResult.Point
	}
	return totalPoint
}
