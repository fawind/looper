package cmd

import (
	"github.com/spf13/cobra"
	"time"
)

// MITMProxy holds the name of the docker-service for recording requests
const MITMProxy = "mitm-proxy"

// Execute executes the root command for the CLI app
func Execute() {
	var (
		service        string
		testCmd        string
		port           int
		outFile        string
		serviceCompose string
		dockerSleep    int
	)
	var rootCmd = &cobra.Command{Use: "app"}
	var cmdRecord = &cobra.Command{
		Use:   "record",
		Short: "Run in record mode",
		Run: func(cmd *cobra.Command, args []string) {
			runRecord(service, testCmd, port, outFile, serviceCompose, dockerSleep)
		},
	}
	var cmdReplay = &cobra.Command{
		Use:   "replay",
		Short: "Run in replay mode",
		Run: func(cmd *cobra.Command, args []string) {
			runReplay(service, testCmd, port, outFile, serviceCompose, dockerSleep)
		},
	}

	rootCmd.AddCommand(cmdRecord)
	rootCmd.AddCommand(cmdReplay)

	rootCmd.PersistentFlags().StringVar(
		&service, "service", "", "Docker service name to test (required)")
	rootCmd.PersistentFlags().StringVar(
		&testCmd, "test", "", "Test command to execute (required)")
	rootCmd.PersistentFlags().IntVar(
		&port, "port", 9999, "Port to use for the MITM proxy")
	rootCmd.PersistentFlags().StringVar(
		&outFile, "out", "out.mitmdump", "File name for the mitm output file")
	rootCmd.PersistentFlags().StringVar(
		&serviceCompose, "compose", "./docker-compose.yml", "Default docker-compose file for the services")
	rootCmd.PersistentFlags().IntVar(
		&dockerSleep, "sleep", 0, "Time to wait after starting docker services in ms")
	rootCmd.MarkPersistentFlagRequired("service")
	rootCmd.MarkPersistentFlagRequired("test")

	rootCmd.Execute()
}

func runRecord(service string, testCmd string, port int, outFile string, serviceCompose string, dockerSleep int) {
	proxyComposeContent := GetRecordCompose(service, port, outFile)
	ExecTest(serviceCompose, proxyComposeContent, testCmd, time.Duration(dockerSleep)*time.Millisecond, []string{})
}

func runReplay(service string, testCmd string, port int, outFile string, serviceCompose string, dockerSleep int) {
	proxyComposeContent := GetReplayCompose(service, port, outFile)
	ExecTest(serviceCompose, proxyComposeContent, testCmd, time.Duration(dockerSleep)*time.Millisecond, []string{service, MITMProxy})
}
