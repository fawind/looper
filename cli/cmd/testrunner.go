package cmd

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ExecTest starts the docker images and executes the tests
func ExecTest(serviceCompose string, proxyComposeContent string, testCmdInput string, dockerSleep time.Duration) {
	proxyFile := createTmpFile(proxyComposeContent)
	defer os.Remove(proxyFile.Name())
	dockerCmd := startDocker(serviceCompose, proxyFile.Name())
	if dockerSleep.Nanoseconds() > 0 {
		log.Println("Waiting for docker services")
		time.Sleep(dockerSleep)
	}
	testCmd := executeTestCmd(testCmdInput)
	testErr := testCmd.Wait()
	log.Println("Test command finished. Stopping docker services...")
	if err := dockerCmd.Process.Signal(os.Interrupt); err != nil {
		log.Fatal(errors.Wrap(err, "Error closing docker command"))
	}
	if testErr != nil {
		log.Fatal(errors.Wrap(testErr, "Error executing test command"))
	}
}

func executeTestCmd(testCmdInput string) *exec.Cmd {
	cmd := exec.Command("sh", "-c", testCmdInput)
	execCommand(cmd)
	return cmd
}

func startDocker(serviceCompose string, proxyCompose string) *exec.Cmd {
	cmd := exec.Command("docker-compose", "-f", serviceCompose, "-f", proxyCompose, "up")
	execCommand(cmd)
	return cmd
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
	log.Printf("Running command: \"%s\"\n", strings.Join(cmd.Args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(errors.Wrap(err, "Error executing command"))
	}
}
