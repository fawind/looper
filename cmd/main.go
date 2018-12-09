package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func main() {
	var (
		service string
		port    int
		outFile string
	)
	var rootCmd = &cobra.Command{Use: "app"}
	var cmdRecord = &cobra.Command{
		Use:   "record",
		Short: "Run in record mode",
		Run: func(cmd *cobra.Command, args []string) {
			runRecord(service, port, outFile)
		},
	}
	var cmdReplay = &cobra.Command{
		Use:   "replay",
		Short: "Run in replay mode",
		Run: func(cmd *cobra.Command, args []string) {
			runReplay(service, port, outFile)
		},
	}

	rootCmd.AddCommand(cmdRecord)
	rootCmd.AddCommand(cmdReplay)

	rootCmd.PersistentFlags().StringVar(&service, "service", "", "Docker service name to test (required)")
	rootCmd.PersistentFlags().IntVar(&port, "port", 9999, "Port to use for the MITM proxy")
	rootCmd.PersistentFlags().StringVar(&outFile, "out", "out.mitmdump", "File name for the mitm output file")
	rootCmd.MarkPersistentFlagRequired("service")

	rootCmd.Execute()
}

func runRecord(service string, port int, outFile string) {
	var composeFile = GetRecordCompose(service, port, outFile)
	fmt.Printf("Docker Compose: \n%s", composeFile)
}

func runReplay(service string, port int, outFile string) {
	var composeFile = GetReplayCompose(service, port, outFile)
	fmt.Printf("Docker Compose: \n%s", composeFile)
}
