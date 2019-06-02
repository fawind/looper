package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"time"
)

// MITMProxy holds the name of the docker-service for recording requests
const MITMProxy = "mitm-proxy"

// Config contains the config parameters
type Config struct {
	service        string
	testCmd        string
	port           int
	outFile        string
	dumpDir        string
	serviceCompose string
	dockerSleep    int
}

// Execute executes the root command for the CLI app
func Execute() {
	var cfg Config
	var rootCmd = &cobra.Command{Use: "app"}
	var cmdRecord = &cobra.Command{
		Use:   "record",
		Short: "Run in record mode",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.outFile = getDumpFile(cfg.outFile, cfg.service, cfg.testCmd)
			runRecord(cfg)
		},
	}
	var cmdReplay = &cobra.Command{
		Use:   "replay",
		Short: "Run in replay mode",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.outFile = getDumpFile(cfg.outFile, cfg.service, cfg.testCmd)
			runReplay(cfg)
		},
	}

	rootCmd.AddCommand(cmdRecord)
	rootCmd.AddCommand(cmdReplay)

	rootCmd.PersistentFlags().StringVar(
		&cfg.service, "service", "", "Docker service name to test (required)")
	rootCmd.PersistentFlags().StringVar(
		&cfg.testCmd, "test", "", "Test command to execute (required)")
	rootCmd.PersistentFlags().IntVar(
		&cfg.port, "port", 9999, "Port to use for the MITM proxy")
	rootCmd.PersistentFlags().StringVar(
		&cfg.serviceCompose, "compose", "./docker-compose.yml", "Default docker-compose file for the services")
	rootCmd.PersistentFlags().StringVar(
		&cfg.outFile, "out", "", "File name for the mitm output file")
	rootCmd.PersistentFlags().StringVar(
		&cfg.dumpDir, "directory", "replay-dumps", "Directory to store the test dumps in")
	rootCmd.PersistentFlags().IntVar(
		&cfg.dockerSleep, "sleep", 0, "Time to wait after starting docker services in ms")

	rootCmd.MarkPersistentFlagRequired("service")
	rootCmd.MarkPersistentFlagRequired("test")

	rootCmd.Execute()
}

func runRecord(cfg Config) {
	proxyComposeContent := GetRecordCompose(cfg.service, cfg.port, cfg.outFile, cfg.dumpDir)
	ExecTest(
		cfg.service,
		cfg.serviceCompose,
		proxyComposeContent,
		cfg.testCmd,
		time.Duration(cfg.dockerSleep)*time.Millisecond,
		[]string{})
	log.Printf("Request log saved to \"%s/%s\"", cfg.dumpDir, cfg.outFile)
}

func runReplay(cfg Config) {
	log.Printf("Request log saved to \"%s/%s\"", cfg.dumpDir, cfg.outFile)
	if !DumpExists(cfg.dumpDir, cfg.outFile) {
		log.Fatal("Dump file not found. Run record first or manually provide a dump file.")
	}
	proxyComposeContent := GetReplayCompose(cfg.service, cfg.port, cfg.outFile, cfg.dumpDir)
	ExecTest(
		cfg.service,
		cfg.serviceCompose,
		proxyComposeContent,
		cfg.testCmd,
		time.Duration(cfg.dockerSleep)*time.Millisecond,
		[]string{cfg.service, MITMProxy})
}

func getDumpFile(outFile string, service string, testCmd string) string {
	if len(outFile) > 0 {
		return outFile
	}
	return GetDumpFileName(service, testCmd)
}
