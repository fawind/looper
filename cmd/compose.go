package cmd

import (
	"bytes"
	"github.com/pkg/errors"
	"text/template"
)

type composeOptions struct {
	ProxyMode bool
	Service   string
	ProxyName string
	MITMArgs  string
	OutFile   string
	Port      int
	DumpDir   string
}

const composeTemplate string = `version: '3'
services:
  {{.ProxyName}}:
    container_name: "{{.ProxyName}}"
    image: mitmproxy/mitmproxy
    entrypoint: "mitmdump {{.MITMArgs}} /{{.DumpDir}}/{{.OutFile}} -p {{.Port}}"
    {{if .ProxyMode}}network_mode: host{{end}}
    ports:
      - "{{.Port}}:{{.Port}}"
    volumes:
      - ./{{.DumpDir}}:/{{.DumpDir}}
  {{if not .ProxyMode}}{{.Service}}:
    container_name: {{.Service}}
    depends_on:
      - mitm-proxy
    environment:
      http_proxy: http://mitm-proxy:{{.Port}}
  {{end}}
`

type RecordOption func(options *composeOptions)

type ModeOption func(options *composeOptions)

func SetDockerMode(service string) ModeOption {
	return func(args *composeOptions) {
		args.ProxyMode = false
		args.Service = service
	}
}

func SetProxyMode() ModeOption {
	return func(args *composeOptions) {
		args.ProxyMode = true
		args.Service = ""
	}
}

func SetRecord() RecordOption {
	return func(args *composeOptions) {
		args.MITMArgs = getMitmArgs(true)
	}
}

func SetReplay() RecordOption {
	return func(args *composeOptions) {
		args.MITMArgs = getMitmArgs(false)
	}
}

// GetCompose returns the generated docker-compose
func GetCompose(recordSetter RecordOption, modeSetter ModeOption, port int, outFile string, dumpDir string) string {
	var tmplOptions = &composeOptions{OutFile: outFile, Port: port, DumpDir: dumpDir, ProxyName: MITMProxy}
	recordSetter(tmplOptions)
	modeSetter(tmplOptions)
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
