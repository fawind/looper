package cmd

import (
	"github.com/spf13/cobra"
)

// Execute executes the root command for the CLI app
func Execute() {
	var (
		service        string
		port           int
		outFile        string
		serviceCompose string
	)
	var rootCmd = &cobra.Command{Use: "app"}
	var cmdRecord = &cobra.Command{
		Use:   "record",
		Short: "Run in record mode",
		Run: func(cmd *cobra.Command, args []string) {
			runRecord(service, port, outFile, serviceCompose)
		},
	}
	var cmdReplay = &cobra.Command{
		Use:   "replay",
		Short: "Run in replay mode",
		Run: func(cmd *cobra.Command, args []string) {
			runReplay(service, port, outFile, serviceCompose)
		},
	}

	rootCmd.AddCommand(cmdRecord)
	rootCmd.AddCommand(cmdReplay)

	rootCmd.PersistentFlags().StringVar(
		&service, "service", "", "Docker service name to test (required)")
	rootCmd.PersistentFlags().IntVar(
		&port, "port", 9999, "Port to use for the MITM proxy")
	rootCmd.PersistentFlags().StringVar(
		&outFile, "out", "out.mitmdump", "File name for the mitm output file")
	rootCmd.PersistentFlags().StringVar(
		&serviceCompose, "compose", "./docker-compose.yml", "Default docker-compose file for the services")
	rootCmd.MarkPersistentFlagRequired("service")

	rootCmd.Execute()
}

func runRecord(service string, port int, outFile string, serviceCompose string) {
	proxyComposeContent := GetRecordCompose(service, port, outFile)
	ExecTest(serviceCompose, proxyComposeContent)
}

func runReplay(service string, port int, outFile string, serviceCompose string) {
	proxyComposeContent := GetReplayCompose(service, port, outFile)
	ExecTest(serviceCompose, proxyComposeContent)
}
