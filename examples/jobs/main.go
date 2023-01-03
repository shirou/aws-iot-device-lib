// SPDX-License-Identifier: Apache-2.0
package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{Name: "key", Required: true},
		&cli.StringFlag{Name: "cert", Required: true},
		&cli.StringFlag{Name: "ca_file", Required: true},
		&cli.StringFlag{Name: "thing_name", Value: "", Usage: ""},
		&cli.StringFlag{Name: "endpoint", Value: "", Required: true},
		&cli.IntFlag{Name: "port", Value: 8883, Usage: "port number"},
	}

	app := &cli.App{
		Flags: flags,
		Commands: []*cli.Command{
			{
				Name:   "DescribeJobExecution",
				Action: DescribeJobExecution,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "jobid", Required: true},
				},
			},
			{
				Name:   "GetPendingJobExecutions",
				Action: GetPendingJobExecutions,
			},
			{
				Name:   "StartNextPendingJobExecution",
				Action: StartNextPendingJobExecution,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "jobid", Required: true},
				},
			},
			{
				Name:   "UpdateJobExecution",
				Action: UpdateJobExecution,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "jobid", Required: true},
					&cli.StringFlag{Name: "status", Value: "SUCCEEDED"},
				},
			},
			{
				Name:   "JobExecutionsChanged",
				Action: JobExecutionsChanged,
			},
			{
				Name:   "NextJobExecutionChanged",
				Action: NextJobExecutionChanged,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
