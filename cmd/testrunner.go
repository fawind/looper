package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ExecTest starts the docker images and executes the tests
func ExecTest(serviceCompose string, proxyComposeContent string, testCmdInput string, dockerSleep time.Duration, servicesToStart []string) {
	proxyFile := createTmpFile(proxyComposeContent)
	defer os.Remove(proxyFile.Name())
	dockerCmd := startDocker(serviceCompose, proxyFile.Name(), servicesToStart)

	waitForService()

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

// waitForService waits for the docker container under test to report as running
func waitForService() {
	for executeCheckRunningCmd() {
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Docker container reported as running.")
}

func executeCheckRunningCmd() bool {
	cmd := exec.Command("sh", "-c", "docker inspect -f {{.State.Running}} service")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	b, err := strconv.ParseBool(strings.TrimSuffix(string(out), "\n"))
	if err == nil {
		return b
	}
	return false
}

func executeTestCmd(testCmdInput string) *exec.Cmd {
	cmd := exec.Command("sh", "-c", testCmdInput)
	execCommand(cmd)
	return cmd
}

func startDocker(serviceCompose string, proxyCompose string, servicesToStart []string) *exec.Cmd {
	args := []string{"-f", serviceCompose, "-f", proxyCompose, "up"}
	args = append(args, servicesToStart...)

	cmd := exec.Command("docker-compose", args...)
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
