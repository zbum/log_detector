package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dooray-go/dooray"
	"github.com/hpcloud/tail"
	"log"
	"strings"
	"time"
)

func main() {
	filePath := flag.String("f", "", "file path")
	pattern := flag.String("p", "", "pattern")
	hookUrl := flag.String("h", "", "hook url")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("you must specify file path.")
		return
	}

	if *pattern == "" {
		fmt.Println("you must specify pattern.")
		return
	}

	if *hookUrl == "" {
		fmt.Println("you must hook url.")
		return
	}

	t, _ := tail.TailFile(*filePath, tail.Config{
		Location: &tail.SeekInfo{Offset: 0, Whence: 2},
		Follow:   true,
		ReOpen:   true})
	for line := range t.Lines {
		if strings.Contains(line.Text, *pattern) {

			ctx1 := context.Background()
			subCtx1, _ := context.WithTimeout(ctx1, 3*time.Second)
			slackErr := dooray.PostWebhookContext(subCtx1, *hookUrl, &dooray.WebhookMessage{
				BotName: "[log-detector] ",
				Text:    fmt.Sprintf("[DETECT PATTERN: %s] %s", *pattern, line.Text),
			})

			if slackErr != nil {
				log.Printf("dial error: %s", slackErr.Error())
			}

		}
	}
}
