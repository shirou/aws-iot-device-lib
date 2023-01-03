// SPDX-License-Identifier: Apache-2.0
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iotjobsdataplane"
	"github.com/aws/aws-sdk-go-v2/service/iotjobsdataplane/types"
	"github.com/shirou/aws-iot-device-lib/examples/connect"
	"github.com/shirou/aws-iot-device-lib/jobs"
	"github.com/urfave/cli/v2"
)

func getJobsClient(cCtx *cli.Context) (*jobs.Client, error) {
	args := connect.ConnectionArgs{
		Key:       cCtx.String("key"),
		Cert:      cCtx.String("cert"),
		Endpoint:  cCtx.String("endpoint"),
		CAFile:    cCtx.String("ca_file"),
		ThingName: cCtx.String("thing_name"),
		Port:      cCtx.Int("port"),
	}

	mc, err := connect.Connect(args)
	if err != nil {
		return nil, err
	}

	client, err := jobs.NewClient(mc)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func DescribeJobExecution(cCtx *cli.Context) error {
	fmt.Println("--------------------")
	fmt.Println("DescribeJobExecution")
	client, err := getJobsClient(cCtx)
	if err != nil {
		return err
	}
	thingName := cCtx.String("thing_name")
	jobId := cCtx.String("jobid")

	req := jobs.DescribeJobExecutionInput{}

	ctx := context.Background()
	ret, err := client.DescribeJobExecution(ctx, thingName, jobId, req)
	if err != nil {
		return err
	}
	for _, step := range ret.Execution.JobDocument.Steps {
		fmt.Printf("Steps: %s\n", step.Action.Name)
	}
	return nil
}

func GetPendingJobExecutions(cCtx *cli.Context) error {
	fmt.Println("--------------------")
	fmt.Println("GetPendingJobExecutions")
	client, err := getJobsClient(cCtx)
	if err != nil {
		return err
	}

	thingName := cCtx.String("thing_name")
	req := iotjobsdataplane.GetPendingJobExecutionsInput{}

	ctx := context.Background()
	ret, err := client.GetPendingJobExecutions(ctx, thingName, req)
	if err != nil {
		return err
	}
	if len(ret.QueuedJobs) == 0 {
		fmt.Println("No queued jobs")
	}
	for _, r := range ret.QueuedJobs {
		fmt.Printf("PendingJobs: JobID=%s, QueuedAt=%d\n", *r.JobId, r.QueuedAt)
	}
	return nil
}

func StartNextPendingJobExecution(cCtx *cli.Context) error {
	fmt.Println("--------------------")
	fmt.Println("StartNextPendingJobExecution")
	client, err := getJobsClient(cCtx)
	if err != nil {
		return err
	}
	thingName := cCtx.String("thing_name")

	req := iotjobsdataplane.StartNextPendingJobExecutionInput{}

	ctx := context.Background()
	ret, err := client.StartNextPendingJobExecution(ctx, thingName, req)
	if err != nil {
		return err
	}
	fmt.Println(ret)
	return nil
}

func UpdateJobExecution(cCtx *cli.Context) error {
	fmt.Println("--------------------")
	fmt.Println("UpdateJobExecution")

	client, err := getJobsClient(cCtx)
	if err != nil {
		return err
	}
	req := jobs.UpdateJobExecutionInput{
		Status: types.JobExecutionStatus(cCtx.String("status")),
	}
	ctx := context.Background()
	ret, err := client.UpdateJobExecution(ctx, cCtx.String("thing_name"), cCtx.String("jobid"), req)
	if err != nil {
		return err
	}
	fmt.Println(ret)
	return nil
}

func JobExecutionsChanged(cCtx *cli.Context) error {
	client, err := getJobsClient(cCtx)
	if err != nil {
		return err
	}

	thingName := cCtx.String("thing_name")

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()
	callback := func(jcli *jobs.Client, msg jobs.JobExecutionsChangedMessage) error {
		ctx := context.Background()
		for _, job := range msg.Jobs[jobs.JobExecutionStatusQueued] {
			req := jobs.DescribeJobExecutionInput{}

			fmt.Println("--------------------")
			fmt.Println("DescribeJobExecution")
			j, _ := jcli.DescribeJobExecution(ctx, thingName, *job.JobId, req)
			for _, step := range j.Execution.JobDocument.Steps {
				fmt.Printf("step: %s\n", step.Action.Name)
			}

			fmt.Println("--------------------")
			fmt.Println("UpdateJobExecution")
			updateReq := jobs.UpdateJobExecutionInput{
				Status: types.JobExecutionStatus(jobs.JobExecutionStatusSucceeded),
			}
			if _, err := jcli.UpdateJobExecution(ctx, thingName, *job.JobId, updateReq); err != nil {
				fmt.Println(err)
			}
			fmt.Println("success")
		}
		return nil
	}
	go client.JobExecutionsChanged(ctx, thingName, callback)

	<-ctx.Done()
	return ctx.Err()
}

func NextJobExecutionChanged(cCtx *cli.Context) error {
	client, err := getJobsClient(cCtx)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	callback := func(client *jobs.Client, msg jobs.NextJobExecutionChangedMessage) error {
		return nil
	}
	go client.NextJobExecutionChanged(ctx, cCtx.String("thing_name"), callback)

	<-ctx.Done()
	return ctx.Err()
}
