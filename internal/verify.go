package internal

func Verifiy(file string, options *Options) error {
	tm, error := NewTaskManager(file, options)
	if error != nil {
		return error
	}
	tm.Process()

	return nil
}
