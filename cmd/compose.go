package cmd

import (
	"bytes"
	"github.com/pkg/errors"
	"text/template"
)

type composeOptions struct {
	Service string
	MitmArg string
	OutFile string
	Port    int
	DumpDir string
}

const composeTemplate string = `version: '3'
services:
  mitm-proxy:
    image: mitmproxy/mitmproxy
    entrypoint: "mitmdump {{.MitmArg}} /{{.DumpDir}}/{{.OutFile}} -p {{.Port}}"
    ports:
      - "{{.Port}}:{{.Port}}"
    volumes:
      - ./{{.DumpDir}}:/{{.DumpDir}}
  {{.Service}}:
    container_name: {{.Service}}
    depends_on:
      - mitm-proxy
    environment:
      http_proxy: http://mitm-proxy:{{.Port}}
`

// GetRecordCompose returns the docker-compose config for record mode
func GetRecordCompose(service string, port int, outFile string, dumpDir string) string {
	return getCompose(true, service, port, outFile, dumpDir)
}

// GetReplayCompose returns the docker-compose config for replaying mode
func GetReplayCompose(service string, port int, outFile string, dumpDir string) string {
	return getCompose(false, service, port, outFile, dumpDir)
}

func getCompose(isRecord bool, service string, port int, outFile string, dumpDir string) string {
	var tmplOptions = composeOptions{service, getMitmArgs(isRecord), outFile, port, dumpDir}
	var tmpl, err = template.New("MITM Docker Compose").Parse(composeTemplate)
	if err != nil {
		panic(errors.Wrap(err, "Error parsing template"))
	}
	var out bytes.Buffer
	if err := tmpl.Execute(&out, tmplOptions); err != nil {
		panic(errors.Wrap(err, "Error executing template"))
	}
	return out.String()
}

func getMitmArgs(isRecord bool) string {
	const (
		record = "-w"
		replay = "-S"
	)
	if isRecord {
		return record
	}
	return replay
}
