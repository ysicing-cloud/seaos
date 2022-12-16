package k3s

type Server struct {
	metadata MetaData
}

func (s Server) Template() string {
	return render(s.metadata, master)
}
