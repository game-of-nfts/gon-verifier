package scorecard

import (
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	// entrance directory
	DefaultScoreCardFile = "scorecard.xlsx"
	DefaultEvidenceFile  = "evidence.xlsx"
	DefaultRankIndivOne  = "rankB3.xlsx"
	DefaultRankIndivTwo  = "rankB4.xlsx"
	DefaultRankTeamOne   = "rankB8.xlsx"
	DefaultQuizGame      = "rankB9.xlsx"
	// user directory
	DefaultStageOneTaskPoint   = "taskpoint1.xlsx"
	DefaultStageTwoTaskPoint   = "taskpoint2.xlsx"
	DefaultStageTwoBTaskPoint  = "taskpoint2b.xlsx"
	DefaultStageThreeTaskPoint = "taskpoint3.xlsx"
	// sheet name
	DefaultTaskPointSheet = "result"
	DefaultScoreCardSheet = "result"
)

// ScoreCard reads task results and output to the scorecard
type ScoreCard struct {
	entranceDir string // path: entrance
}

type ScoreCardEntry struct {
	taskCompleted string
	totalPoint    int
	teamName      string
	failedReason  string
	githubAccount string
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
	f.SetCellValue(DefaultScoreCardSheet, "F1", "failed_reason")
	f.SetCellValue(DefaultScoreCardSheet, "G1", "github_account")

	var allTaskPointFiles [][]string
	var taskPointFiles []string
	err = filepath.Walk(sc.entranceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == DefaultStageOneTaskPoint {
			taskPointFiles = append(taskPointFiles, path)
		}

		if !info.IsDir() && info.Name() == DefaultStageTwoTaskPoint {
			taskPointFiles = append(taskPointFiles, path)
		}

		if !info.IsDir() && info.Name() == DefaultStageTwoBTaskPoint {
			taskPointFiles = append(taskPointFiles, path)
		}

		if info.IsDir() && len(taskPointFiles) != 0 {
			allTaskPointFiles = append(allTaskPointFiles, taskPointFiles)
			taskPointFiles = []string{}
		}
		return nil
	})
	allTaskPointFiles = append(allTaskPointFiles, taskPointFiles)
	if err != nil {
		return err
	}

	var entries []*ScoreCardEntry
	for _, taskPointFiles := range allTaskPointFiles {
		entry, err := sc.HandleTaskPoint(taskPointFiles)
		if err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].totalPoint != entries[j].totalPoint {
			return entries[i].totalPoint > entries[j].totalPoint
		}
		return entries[i].teamName < entries[j].teamName
	})

	rank := 1
	for i := 0; i < len(entries); i++ {
		if i != 0 {
			if entries[i].totalPoint != entries[i-1].totalPoint {
				rank++
			}
		}

		f.SetCellValue(DefaultScoreCardSheet, "A"+strconv.Itoa(i+2), rank)
		f.SetCellValue(DefaultScoreCardSheet, "B"+strconv.Itoa(i+2), entries[i].teamName)
		f.SetCellValue(DefaultScoreCardSheet, "C"+strconv.Itoa(i+2), entries[i].taskCompleted)
		f.SetCellValue(DefaultScoreCardSheet, "D"+strconv.Itoa(i+2), entries[i].totalPoint)
		f.SetCellValue(DefaultScoreCardSheet, "F"+strconv.Itoa(i+2), entries[i].failedReason)
		f.SetCellValue(DefaultScoreCardSheet, "G"+strconv.Itoa(i+2), entries[i].githubAccount)
	}

	f.SetActiveSheet(index)
	scorecardFile := filepath.Join(sc.entranceDir, DefaultScoreCardFile)
	if err := f.SaveAs(scorecardFile); err != nil {
		return err
	}

	return nil
}

// HandleTaskPoint calculate the score of a stage from one's task point files
func (sc *ScoreCard) HandleTaskPoint(taskPointFiles []string) (*ScoreCardEntry, error) {
	var (
		taskResults TaskResults
		github      string
		teamName    string
	)

	for i, taskPointFile := range taskPointFiles {
		taskpoint, err := excelize.OpenFile(taskPointFile)
		if err != nil {
			return nil, err
		}
		defer taskpoint.Close()

		rows, err := taskpoint.GetRows(DefaultTaskPointSheet)
		if err != nil {
			return nil, err
		}

		strs := strings.Split(taskPointFile, "/")
		github = strs[len(strs)-2]

		if i == 0 && len(rows) == 1 {
			return &ScoreCardEntry{
				taskCompleted: "",
				totalPoint:    0,
				teamName:      "UnknownTeam:@" + github,
				failedReason:  "all evidence formats are incorrect",
				githubAccount: "@" + github,
			}, nil
		}

		if len(rows) == 1 {
			continue
		}

		for _, row := range rows[1:] {
			point, _ := strconv.Atoi(row[2])
			reason := ""
			if len(row) == 4 {
				reason = row[3]
			}
			taskResults = append(taskResults, TaskResult{
				TaskNo: row[0],
				Point:  point,
				Reason: reason,
			})
		}

		teamName = rows[1][1]
		if strings.HasPrefix(teamName, "team") {
			teamName = "UnknownTeam:@" + github
		}
	}

	return &ScoreCardEntry{
		taskCompleted: sc.concatenateTaskNo(taskResults),
		totalPoint:    sc.calculateTotalPoint(taskResults),
		teamName:      teamName,
		failedReason:  sc.concatenateFailedReason(taskResults),
		githubAccount: "@" + github,
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

func (sc *ScoreCard) concatenateFailedReason(taskResults TaskResults) string {
	// NOTE: call after concatenateTaskNo
	failedReasons := make([]string, 0)
	reasons := make(map[string]bool, 0)
	for i := 0; i < len(taskResults); i++ {
		if len(taskResults[i].Reason) == 0 {
			continue
		}
		failedReason := taskResults[i].TaskNo + ":" + taskResults[i].Reason
		reasons[failedReason] = true
	}
	for key := range reasons {
		failedReasons = append(failedReasons, key)
	}

	return strings.Join(failedReasons, ",")
}

func (sc *ScoreCard) calculateTotalPoint(taskResults TaskResults) int {
	var totalPoint int
	for _, taskResult := range taskResults {
		totalPoint += taskResult.Point
	}
	return totalPoint
}
