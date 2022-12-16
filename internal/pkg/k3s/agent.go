package k3s

type Agent struct {
	metadata MetaData
}

func (a Agent) Template() string {
	return render(a.metadata, agnet)
}
