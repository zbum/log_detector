package main

import (
	"context"
	"fmt"
	"github.com/cosiner/flag"
	"github.com/dooray-go/dooray"
	"github.com/hpcloud/tail"
	"log"
	"os"
	"strings"
	"time"
)

type Flags struct {
	FilePath string `names:"-f, --filepath" important:"1" usage:"log file path"`
	Pattern  string `names:"-p, --pattern" usage:"detection pattern"`
	HookUrl  string `names:"-u, --hookurl" usage:"dooray hook url"`
}

func (t *Flags) Metadata() map[string]flag.Flag {
	const (
		usage   = "log_detector is a tool for monitoring log file and alert."
		version = `
			version: v1.0.0
			date:   2024-04-28 10:00:01
		`
		desc = `
		log_detector is a tool for monitoring log file and alert. This use Dooray Messenger
		`
	)
	return map[string]flag.Flag{
		"": {
			Usage:   usage,
			Version: version,
			Desc:    desc,
		},
	}
}

func main() {
	var flags Flags
	flagSet := flag.NewFlagSet(flag.Flag{})
	err := flagSet.ParseStruct(&flags, os.Args...)
	if err != nil {
		flagSet.Help()
		return
	}
	if flags.HookUrl == "" || flags.Pattern == "" || flags.FilePath == "" {
		flagSet.Help()
		return
	}

	t, _ := tail.TailFile(flags.FilePath, tail.Config{
		Location: &tail.SeekInfo{Offset: 0, Whence: 2},
		Follow:   true,
		ReOpen:   true})
	for line := range t.Lines {
		if strings.Contains(line.Text, flags.Pattern) {

			ctx1 := context.Background()
			subCtx1, _ := context.WithTimeout(ctx1, 3*time.Second)
			doorayErr := dooray.PostWebhookContext(subCtx1, flags.HookUrl, &dooray.WebhookMessage{
				BotName: "[log_detector] ",
				Text:    fmt.Sprintf("[DETECT PATTERN: %s] %s", flags.Pattern, line.Text),
			})

			if doorayErr != nil {
				log.Printf("dial error: %s", doorayErr.Error())
			}

		}
	}
}
