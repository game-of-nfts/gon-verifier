package rank

import (
	"errors"
	"fmt"
	"github.com/taramakage/gon-verifier/internal/chain"
	"github.com/taramakage/gon-verifier/internal/scorecard"
	"github.com/xuri/excelize/v2"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type QuizRanker struct {
	TaskNo    string
	TaskPoint int32
	Entrance  string
	Quizers   []Quizer
	r         *chain.Registry
	f         *chain.Flow
}

type Quizer struct {
	TeamName string
	Address  string
	Count    int
	Path     string
}

func NewQuizRanker(entrance, taskNo string, taskPoint int32) *QuizRanker {
	f, err := chain.NewFlow(chain.FlowStrMap["f04"])
	if err != nil {
		return nil
	}

	return &QuizRanker{
		Entrance:  entrance,
		TaskNo:    taskNo,
		TaskPoint: taskPoint,
		Quizers:   make([]Quizer, 0),
		f:         f,
		r:         chain.NewRegistry(),
	}
}

func (qr *QuizRanker) Do() error {
	err := qr.Collect()
	if err != nil {
		return err
	}
	qr.Sort()

	err = qr.GenerateRank()
	if err != nil {
		return err
	}

	err = qr.WriteTaskPoint()
	if err != nil {
		return err
	}
	return nil
}

func (qr *QuizRanker) Collect() error {
	err := qr.collectNft()
	if err != nil {
		return err
	}
	qr.collectQuizers()
	return nil
}

// collectNftAddr query data on-chain and count the number of NFTs and their ownership information.
func (qr *QuizRanker) collectNft() error {
	hash, err := qr.f.GetFinalIbcHash("gonQuiz")
	if err != nil {
		return err
	}
	ibcClassId := "ibc/" + hash.String()
	c := qr.r.GetChain(chain.ChainIdAbbreviationIris)
	iris, ok := c.(*chain.Iris)
	if !ok {
		return errors.New("failed to get chain")
	}

	collection, err := iris.GetCollection(ibcClassId)
	if err != nil {
		return err
	}

	nfts := collection.Collection.NFTs
	addrMap := make(map[string]int)
	for _, nft := range nfts {
		addr := nft.Owner
		addrMap[addr] += 1
	}

	for addr, count := range addrMap {
		qr.Quizers = append(qr.Quizers, Quizer{
			Address: addr,
			Count:   count,
		})
	}

	return nil
}

// collectQuizer walk participants' evidence file and read info of quizer
func (qr *QuizRanker) collectQuizers() {
	strs := make([]string, 0)

	err := filepath.Walk(qr.Entrance, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() || info.Name() != scorecard.DefaultEvidenceFile {
			return nil
		}
		strs = append(strs, path)
		return nil
	})
	if err != nil {
		err.Error()
	}

	for _, str := range strs {
		qr.collectQuizer(str)
	}
}

func (qr *QuizRanker) collectQuizer(evidence string) {
	f, err := excelize.OpenFile(evidence)
	defer f.Close()

	rows, err := f.GetRows("Info")
	if err != nil || len(rows) < 2 {
		return
	}

	teamName := rows[1][0]
	addr := rows[1][1]

	for i, quizer := range qr.Quizers {
		if len(quizer.TeamName) != 0 {
			continue
		}
		if quizer.Address == addr {
			qr.Quizers[i].TeamName = teamName
			qr.Quizers[i].Path = evidence
			break
		}
	}
}

func (qr *QuizRanker) Sort() {
	sort.SliceStable(qr.Quizers, func(i, j int) bool {
		return qr.Quizers[i].Count > qr.Quizers[j].Count
	})
}

func (qr *QuizRanker) GenerateRank() error {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "result"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	f.SetCellValue(sheetName, "A1", "Rank")
	f.SetCellValue(sheetName, "B1", "TeamName")
	f.SetCellValue(sheetName, "C1", "QuizContent")

	for i, quizer := range qr.Quizers {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), quizer.TeamName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), quizer.Count)
	}

	f.SetActiveSheet(index)

	fileName := filepath.Join(qr.Entrance, fmt.Sprintf("rank%s.xlsx", qr.TaskNo))
	if err := f.SaveAs(fileName); err != nil {
		return err
	}
	return nil
}

func (qr *QuizRanker) WriteTaskPoint() error {
	for _, quizer := range qr.Quizers {
		err := qr.clearLegacyTaskPoint(quizer)
		if err != nil {
			return err
		}
		err = qr.writeTaskPoint(quizer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (qr *QuizRanker) clearLegacyTaskPoint(quizer Quizer) error {
	dir := path.Dir(quizer.Path)
	taskpoint3 := path.Join(dir, "taskpoint3.xlsx")
	file, err := excelize.OpenFile(taskpoint3)
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
		if strings.HasPrefix(rows[i][0], qr.TaskNo) {
			file.RemoveRow(sheetName, i+1)
		}
	}

	err = file.Save()
	if err != nil {
		return err
	}

	return nil
}

func (qr *QuizRanker) writeTaskPoint(quizer Quizer) error {
	dir := path.Dir(quizer.Path)
	taskpoint3 := path.Join(dir, "taskpoint3.xlsx")
	file, err := excelize.OpenFile(taskpoint3)
	if err != nil {
		return err
	}
	defer file.Close()

	sheetName := "result"
	rows, err := file.GetRows(sheetName)
	length := len(rows)

	file.SetCellValue(sheetName, fmt.Sprintf("A%d", length+1), qr.TaskNo+"*"+strconv.Itoa(quizer.Count))
	file.SetCellValue(sheetName, fmt.Sprintf("B%d", length+1), quizer.TeamName)
	file.SetCellValue(sheetName, fmt.Sprintf("C%d", length+1), qr.TaskPoint*int32(quizer.Count))

	err = file.Save()
	if err != nil {
		return err
	}

	return nil
}
