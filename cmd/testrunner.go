package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// ExecTest starts the docker images and executes the tests
// TODO(fawind): Add param for test cmd target
func ExecTest(serviceCompose string, proxyComposeContent string) {
	proxyFile := createTmpFile(proxyComposeContent)
	defer os.Remove(proxyFile.Name())
	startDocker(serviceCompose, proxyFile.Name())
}

func startDocker(serviceCompose string, proxyCompose string) {
	cmd := exec.Command("docker-compose", "-f", serviceCompose, "-f", proxyCompose, "up")
	execCommand(cmd)
}

func createTmpFile(content string) *os.File {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "docker-compose-proxy-")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Could not create tempfile"))
	}
	if _, err = tmpFile.Write([]byte(content)); err != nil {
		log.Fatal(errors.Wrap(err, "Could not create tempfile"))
	}
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}
	return tmpFile
}

func execCommand(cmd *exec.Cmd) {
	fmt.Printf("Running command: `%s`\n", strings.Join(cmd.Args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(errors.Wrap(err, "Error executing command"))
	}
}
