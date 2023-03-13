package verifier

import (
	"errors"
	"fmt"
	"github.com/taramakage/gon-verifier/internal/scorecard"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/exp/slog"

	"github.com/taramakage/gon-verifier/internal/chain"
)

type (
	Options struct {
		TaskNos []string
	}

	Task struct {
		taskNo string
		params any
		vf     Verifier
	}

	TaskManager struct {
		tasks    []Task
		user     UserInfo
		cr       *chain.Registry
		vr       *Registry
		wg       *sync.WaitGroup
		baseDir  string
		resultCh chan *Response
		stopCh   chan int
		saveCh   chan int
	}
)

func NewTaskManager(evidenceFile string, opts *Options) (*TaskManager, error) {
	cr := chain.NewRegistry()
	tm := &TaskManager{
		wg:       &sync.WaitGroup{},
		cr:       cr,
		vr:       NewRegistry(cr),
		resultCh: make(chan *Response, 10),
		stopCh:   make(chan int),
		saveCh:   make(chan int),
	}

	if err := tm.loadEvidence(evidenceFile, opts); err != nil {
		return nil, err
	}
	return tm, nil
}

// Process concurrently verify tasks of one participant and write the result to xlsx file.
func (tm *TaskManager) Process() {
	if len(tm.tasks) == 0 {
		slog.Info("no task process")
		return
	}
	slog.Info("start to verify", "TeamName", tm.user.TeamName, " Github", tm.user.Github)
	go tm.receive()
	for _, task := range tm.tasks {
		tm.wg.Add(1)
		go func(task Task) {
			defer tm.wg.Done()
			// slog.Info("verify rule", "TeamName", tm.user.TeamName, "TaskNo", task.taskNo)
			task.vf.Do(*&Request{
				TaskNo: task.taskNo,
				User:   tm.user,
				Params: task.params,
			}, tm.resultCh)
		}(task)
	}
	tm.wg.Wait()
	tm.stop()
	<-tm.saveCh
	return
}

func (tm *TaskManager) receive() {
	f := excelize.NewFile()

	sheetName := "result"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		slog.Error("NewSheet error", err)
		return
	}

	rowIdx := 1
	f.SetCellValue(sheetName, "A1", "TaskNo")
	f.SetCellValue(sheetName, "B1", "TeamName")
	f.SetCellValue(sheetName, "C1", "Point")
	f.SetCellValue(sheetName, "D1", "Reason")

	for {
		select {
		case result := <-tm.resultCh:
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIdx+1), result.TaskNo)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIdx+1), result.TeamName)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIdx+1), result.Point)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIdx+1), result.Reason)
			rowIdx++
		case <-tm.stopCh:
			f.SetActiveSheet(index)

			fileName := filepath.Join(tm.baseDir, scorecard.DefaultTaskPointFile)
			if err := f.SaveAs(fileName); err != nil {
				slog.Error("Save file error", err)
			}

			if err := f.Close(); err != nil {
				slog.Error("close file error", err)
			}
			tm.saveCh <- 1
			return
		}
	}
}

func (tm *TaskManager) stop() {
	// slog.Info("verify finish", "TeamName", tm.user.TeamName)
	tm.stopCh <- 1
}

func (tm *TaskManager) loadEvidence(evidenceFile string, opts *Options) error {
	evidence, err := excelize.OpenFile(evidenceFile)
	if err != nil {
		return err
	}

	tm.baseDir = filepath.Dir(evidenceFile)
	if err := tm.loadUserInfo(evidence); err != nil {
		return err
	}

	return tm.buildTask(evidence, opts)
}

// loadUserInfo loads the user info from the evidence file.
func (tm *TaskManager) loadUserInfo(evidence *excelize.File) error {
	rows, err := evidence.GetRows("Info")
	if err != nil {
		return errors.New("info sheet not found")
	}

	if len(rows) < 2 {
		return errors.New("info sheet format error")
	}

	columns := rows[1]
	github := columns[0]
	if len(tm.baseDir) > 0 {
		paths := strings.Split(tm.baseDir, string(os.PathSeparator))
		github = paths[len(paths)-1]
	}

	tm.user = UserInfo{
		TeamName: columns[0],
		Github:   github,
		Address: map[string]string{
			chain.ChainIdAbbreviationIris:     columns[1],
			chain.ChainIdAbbreviationStars:    columns[2],
			chain.ChainIdAbbreviationJuno:     columns[3],
			chain.ChainIdAbbreviationUptick:   columns[4],
			chain.ChainIdAbbreviationOmniflix: columns[5],
		},
	}
	return nil
}

// buildTask builds the task list from the evidence file.
func (tm *TaskManager) buildTask(evidence *excelize.File, opts *Options) error {
	taskNos := evidence.GetSheetList()
	if len(opts.TaskNos) != 0 {
		taskNos = opts.TaskNos
	}
	for _, taskNo := range taskNos {
		rowsCols, err := evidence.GetRows(taskNo)
		if err != nil {
			return err
		}

		if len(rowsCols) == 0 {
			return errors.New("evidence sheet is empty")
		}

		vf := tm.vr.Get(taskNo)
		params, err := vf.BuildParams(rowsCols[1:])
		if err != nil {
			return err
		}

		tm.tasks = append(tm.tasks, Task{
			taskNo: taskNo,
			params: params,
			vf:     vf,
		})
	}
	return nil
}

func (tm *TaskManager) Close() {
	for _, v := range tm.cr.GetChains() {
		v.Close()
	}
}
