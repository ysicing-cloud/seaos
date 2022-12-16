package k3s

import (
	"bytes"
	"text/template"
)

type MetaData struct {
	TLSSAN      string
	Master0     bool
	Docker      bool
	DataDir     string
	ClusterCidr string
	ServiceCidr string
	Token       string
	Server      string
}

type Service interface {
	Template() string
}

func NewService(t string, metadata MetaData) Service {
	switch t {
	case "server", "master":
		return &Server{metadata: metadata}
	default:
		return &Agent{metadata: metadata}
	}
}

func render(data MetaData, temp string) string {
	var b bytes.Buffer
	t := template.Must(template.New("k3s").Parse(temp))
	_ = t.Execute(&b, &data)
	return b.String()
}
