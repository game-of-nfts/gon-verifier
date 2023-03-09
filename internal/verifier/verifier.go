package verifier

type GonVerifier struct {
	entrance string
	options  *Options
}

func NewGonVerifier(entrance string, options *Options) *GonVerifier {
	return &GonVerifier{
		entrance: entrance,
		options:  options,
	}
}

func (gv *GonVerifier) Verify(file string) error {
	tm, err := NewTaskManager(file, gv.options)
	if err != nil {
		return err
	}

	if tm != nil {
		tm.Process()
	}

	return nil
}

// TODO: skip some verification based on options
// func (gv *GonVerifier) Skip() {}
