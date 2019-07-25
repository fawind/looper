package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"time"
)

// Required flag names
const serviceFlag = "service"
const testCmdFlag = "test"
const outFileFlag = "out"

// MITMProxy holds the name of the docker-service for recording requests
const MITMProxy = "mitm-proxy"

// Config contains the config parameters for record and replay modes
type Config struct {
	proxyOnly      bool
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
			verifyFlags(cfg, cmd)
			cfg.outFile = getDumpFile(cfg.outFile, cfg.service, cfg.testCmd)
			runRecord(cfg)
		},
	}
	var cmdReplay = &cobra.Command{
		Use:   "replay",
		Short: "Run in replay mode",
		Run: func(cmd *cobra.Command, args []string) {
			verifyFlags(cfg, cmd)
			cfg.outFile = getDumpFile(cfg.outFile, cfg.service, cfg.testCmd)
			runReplay(cfg)
		},
	}

	// Flags for record and replay mode
	rootCmd.PersistentFlags().BoolVar(
		&cfg.proxyOnly, "proxyOnly", false, "Only start proxy in stand-alone mode")
	rootCmd.PersistentFlags().StringVar(
		&cfg.service, serviceFlag, "", "Docker setService name to test (required when not run in stand-alone proxy mode)")
	rootCmd.PersistentFlags().StringVar(
		&cfg.testCmd, testCmdFlag, "", "Test command to execute (required when not run in stand-alone proxy mode)")
	rootCmd.PersistentFlags().IntVar(
		&cfg.port, "port", 9999, "Port to use for the MITM proxy")
	rootCmd.PersistentFlags().StringVar(
		&cfg.serviceCompose, "compose", "./docker-compose.yml", "Default docker-compose file for the services")
	rootCmd.PersistentFlags().StringVar(
		&cfg.outFile, outFileFlag, "", "File name for the mitm output file (required when run in stand-alone proxy mode)")
	rootCmd.PersistentFlags().StringVar(
		&cfg.dumpDir, "directory", "replay-dumps", "Directory to store the test dumps in")
	rootCmd.PersistentFlags().IntVar(
		&cfg.dockerSleep, "sleep", 0, "Time to wait after starting docker services in ms")

	rootCmd.AddCommand(cmdRecord, cmdReplay)
	rootCmd.Execute()
}

func runRecord(cfg Config) {
	log.Printf("Request log saved to \"%s/%s\"", cfg.dumpDir, cfg.outFile)
	proxyComposeContent := getComposeForCfg(true, cfg)
	if cfg.proxyOnly {
		log.Printf("Starting stand-alone proxy mode. Use Ctrl-c to close")
		ExecProxy(proxyComposeContent)
	} else {
		ExecTest(
			cfg.service,
			cfg.serviceCompose,
			proxyComposeContent,
			cfg.testCmd,
			time.Duration(cfg.dockerSleep)*time.Millisecond,
			[]string{})
	}
}

func runReplay(cfg Config) {
	log.Printf("Request log saved to \"%s/%s\"", cfg.dumpDir, cfg.outFile)
	if !DumpExists(cfg.dumpDir, cfg.outFile) {
		log.Fatal("Dump file not found. Run record first or manually provide a dump file.")
	}
	proxyComposeContent := getComposeForCfg(false, cfg)
	if cfg.proxyOnly {
		log.Printf("Starting stand-alone proxy mode. Use Ctrl-c to close")
		ExecProxy(proxyComposeContent)
	} else {
		ExecTest(
			cfg.service,
			cfg.serviceCompose,
			proxyComposeContent,
			cfg.testCmd,
			time.Duration(cfg.dockerSleep)*time.Millisecond,
			[]string{cfg.service, MITMProxy})
	}
}

func getComposeForCfg(isRecord bool, cfg Config) string {
	var recordOption RecordOption
	if isRecord {
		recordOption = SetRecord()
	} else {
		recordOption = SetReplay()
	}
	var modeOption ModeOption
	if cfg.proxyOnly {
		modeOption = SetProxyMode()
	} else {
		modeOption = SetDockerMode(cfg.service)
	}
	return GetCompose(recordOption, modeOption, cfg.port, cfg.outFile, cfg.dumpDir)
}

func verifyFlags(cfg Config, cmd *cobra.Command) {
	errorMsg := "Error: Flag(s) %s not set "
	hasError := false
	var missingFlags []string
	var errorSuffix string
	if !cfg.proxyOnly {
		if len(cfg.service) == 0 {
			missingFlags = append(missingFlags, "--"+serviceFlag)
			hasError = true
		}
		if len(cfg.testCmd) == 0 {
			missingFlags = append(missingFlags, "--"+testCmdFlag)
			hasError = true
		}
		if hasError {
			errorSuffix = "which are required when not run in stand-alone proxy mode\n"
		}
	} else {
		if len(cfg.outFile) == 0 {
			missingFlags = append(missingFlags, "--"+outFileFlag)
			errorSuffix = "which are required when run in stand-alone proxy mode\n"
			hasError = true
		}
	}
	if hasError {
		fmt.Printf(errorMsg+errorSuffix, strings.Join(missingFlags, ", "))
		cmd.Help()
		os.Exit(0)
	}
}

func getDumpFile(outFile string, service string, testCmd string) string {
	if len(outFile) > 0 {
		return outFile
	}
	return GetDumpFileName(service, testCmd)
}
