package scorecard

type TaskResult struct {
	TaskNo string
	Point  int
	Reason string
}

type TaskResults []TaskResult

func (t TaskResults) Len() int {
	return len(t)
}

func (t TaskResults) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TaskResults) Less(i, j int) bool {
	ti := t[i].TaskNo
	tj := t[j].TaskNo
	if ti[0] != tj[0] {
		return ti[0] < tj[0]
	}
	return t.naturalLess(ti, tj)
}

func (t TaskResults) naturalLess(str1, str2 string) bool {
	num1, num2 := 0, 0
	for _, c := range str1 {
		if c >= '0' && c <= '9' {
			num1 = num1*10 + int(c-'0')
		}
	}
	for _, c := range str2 {
		if c >= '0' && c <= '9' {
			num2 = num2*10 + int(c-'0')
		}
	}
	if num1 != num2 {
		return num1 < num2
	}
	return str1 < str2
}
