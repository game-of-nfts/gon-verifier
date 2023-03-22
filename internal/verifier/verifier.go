package verifier

type GonVerifier struct {
	entrance string
	options  *Options
}

const (
	DenomIdGonIndivRace1 = "gonIndivRace1"
	DenomIdGonIndivRace2 = "gonIndivRace2"
	DenomIdGonTeamRace1 = "gonTeamRace1"
	DenomIdGonTeamRace2 = "gonTeamRace2"
	DenomIdGonTeamRace3 = "gonTeamRace3"
	DenomIdGonQuiz = "gonQuiz"

	StartBlockHeightIndivRace1 = "473000"
	StartBlockHeightIndivRace2 = "489000"
	StartBlockHeightTeamRace = "568000"
	EndBlockHeightIndivRace = "516223"
	EndBlockHeightGame = "671700"

	LastOwner = "iaa1488wwr235vka7j722hzacpk0plxw33ksqyneuz"
)

func NewGonVerifier(entrance string, options *Options) *GonVerifier {
	return &GonVerifier{
		entrance: entrance,
		options:  options,
	}
}

func (gv *GonVerifier) Verify(file string) error {
	tm, err := NewTaskManager(file, gv.options)
	defer func() {
		if tm != nil {
			tm.Close()
		}
	}()

	if err != nil {
		return err
	}

	if tm != nil {
		tm.Process(gv.options)
	}

	return nil
}
