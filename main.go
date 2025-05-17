package main

import (
	"context"
	"fmt"
	"github.com/dooray-go/dooray"
	"github.com/hpcloud/tail"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"time"
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "log_detector",
		Short: "log_detector is a tool for monitoring log file and alert.",
		Long:  "log_detector is a tool for monitoring log file and alert.",
		RunE:  errorHandler,
		Run:   executeLogDetector,
	}

	rootCmd.PersistentFlags().StringP("pattern", "p", "", "Pattern to search in log file")
	rootCmd.PersistentFlags().StringP("file", "f", "", "LogFile Path log file")
	rootCmd.PersistentFlags().StringP("hookurl", "u", "", "Dooray Hook URL")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}

}

func errorHandler(cmd *cobra.Command, args []string) error {
	pattern, _ := cmd.Flags().GetString("pattern")
	file, _ := cmd.Flags().GetString("file")
	hookurl, _ := cmd.Flags().GetString("hookurl")

	if pattern == "" || file == "" || hookurl == "" {
		return cmd.Help()
	}
	return nil
}

func executeLogDetector(cmd *cobra.Command, args []string) {
	pattern, _ := cmd.Flags().GetString("pattern")
	file, _ := cmd.Flags().GetString("file")
	hookurl, _ := cmd.Flags().GetString("hookurl")

	t, _ := tail.TailFile(file, tail.Config{
		Location: &tail.SeekInfo{Offset: 0, Whence: 2},
		Follow:   true,
		ReOpen:   true})
	for line := range t.Lines {
		if strings.Contains(line.Text, pattern) {

			ctx1 := context.Background()
			subCtx1, _ := context.WithTimeout(ctx1, 3*time.Second)
			doorayErr := dooray.PostWebhookContext(subCtx1, hookurl, &dooray.WebhookMessage{
				BotName: "[log_detector] ",
				Text:    fmt.Sprintf("[DETECT PATTERN: %s] %s", pattern, line.Text),
			})

			if doorayErr != nil {
				log.Printf("dial error: %s", doorayErr.Error())
			}

		}
	}
}
