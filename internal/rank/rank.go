package rank

type Ranker interface {
	Collect() error
	GenerateRank() error
	WriteTaskPoint() error
}

func Rank(ranker Ranker) error {
	err := ranker.Collect()
	if err != nil {
		return err
	}

	ranker.GenerateRank()
	if err != nil {
		return err
	}

	ranker.WriteTaskPoint()
	if err != nil {
		return err
	}
	return nil
}