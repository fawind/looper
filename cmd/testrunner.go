package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const defaultRetries = 600 // 10 minutes

// ExecTest starts the docker images and executes the tests
func ExecTest(service string, serviceCompose string, proxyComposeContent string, testCmdInput string, dockerSleep time.Duration, servicesToStart []string) {
	proxyFile := createTmpFile(proxyComposeContent)
	defer os.Remove(proxyFile.Name())
	dockerCmd := startDocker(serviceCompose, proxyFile.Name(), servicesToStart)

	waitForService(service)

	if dockerSleep.Nanoseconds() > 0 {
		log.Println("Sleep before executing tests")
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
func waitForService(serviceName string) {
	retries := defaultRetries
	log.Println("Waiting for docker services to be started")
	for executeCheckRunningCmd(serviceName, retries) && retries > 0 {
		retries--
		time.Sleep(1 * time.Second)
	}
	if retries == 0 {
		log.Fatal("Error waiting for docker services, timeout reached")
	}
	log.Println("Docker services reported as running")
}

func executeCheckRunningCmd(serviceName string, retries int) bool {
	cmd := exec.Command("sh", "-c", "docker inspect -f {{.State.Running}} "+serviceName)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	isRunning, err := strconv.ParseBool(strings.TrimSuffix(string(out), "\n"))
	if err != nil {
		return false
	}
	return isRunning
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
