// SPDX-License-Identifier: Apache-2.0
package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iotjobsdataplane"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/shirou/aws-iot-device-lib/internal/mqttutils"
)

const defaultTimeout = 1 * time.Second

type Client struct {
	mc       mqtt.Client
	Timeouts time.Duration
}

func NewClient(mc mqtt.Client) (*Client, error) {
	client := &Client{
		mc:       mc,
		Timeouts: defaultTimeout,
	}

	return client, nil
}

// SetTimeout sets the timeout before a response is returned for an accepted or rejected topic.
// The default is 1 second. In a slow connection environment, it is recommended to set a longer time.
func (client *Client) SetTimeout(dur time.Duration) {
	client.Timeouts = dur
}

type outputType interface {
	DescribeJobExecutionOutput |
		iotjobsdataplane.GetPendingJobExecutionsOutput |
		iotjobsdataplane.UpdateJobExecutionOutput |
		StartNextPendingJobExecutionOutput
}

// handleAsync is a generic processing function. It is not recommended to use this function from outside of this "jobs" package. It may be moved under "internal" in the future.
func handleAsync[K outputType](ctx context.Context, mc mqtt.Client, payload []byte, subTopics []string, pubTopic string) (ret K, err error) {
	accepted := make(chan K)
	rejected := make(chan error)
	callback := func(mc mqtt.Client, msg mqtt.Message) {
		if err := IsError(msg.Payload()); err != nil {
			rejected <- err
			return
		}
		var output K
		if err = json.Unmarshal(msg.Payload(), &output); err != nil {
			rejected <- err
			return
		}

		if strings.HasSuffix(msg.Topic(), "accepted") {
			accepted <- output
		} else if strings.HasSuffix(msg.Topic(), "rejected") {
			rejected <- fmt.Errorf("rejected") // TODO: what payload if rejected?
		} else {
			rejected <- fmt.Errorf("unknown topic subscribed, %s", msg.Topic())
			return
		}
	}

	if err = mqttutils.Subscribe(mc, subTopics, 0, callback); err != nil {
		return
	}
	defer func() {
		close(accepted)
		close(rejected)
		err = mqttutils.JoinErrors(err, mqttutils.Unsubscribe(mc, subTopics))
	}()

	if err = mqttutils.Publish(mc, pubTopic, 0, payload); err != nil {
		return
	}
	for {
		select {
		case r := <-accepted:
			return r, nil
		case err = <-rejected:
			return ret, err
		case <-ctx.Done():
			return ret, ctx.Err()
		}
	}
}

// GetPendingJobExecutions gets detailed information about a job execution.
func (client *Client) GetPendingJobExecutions(ctx context.Context, thingName string, req iotjobsdataplane.GetPendingJobExecutionsInput) (ret iotjobsdataplane.GetPendingJobExecutionsOutput, err error) {
	topics := []string{
		fmt.Sprintf("$aws/things/%s/jobs/get/accepted", thingName),
		fmt.Sprintf("$aws/things/%s/jobs/get/rejected", thingName),
	}
	pubTopic := fmt.Sprintf("$aws/things/%s/jobs/get", thingName)

	payload, err := json.Marshal(req)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, client.Timeouts)
	defer cancel()
	return handleAsync[iotjobsdataplane.GetPendingJobExecutionsOutput](ctx, client.mc, payload, topics, pubTopic)
}

// StartNextPendingJobExecution gets and starts the next pending job execution for a thing
func (client *Client) StartNextPendingJobExecution(ctx context.Context, thingName string, req iotjobsdataplane.StartNextPendingJobExecutionInput) (ret StartNextPendingJobExecutionOutput, err error) {
	topics := []string{
		fmt.Sprintf("$aws/things/%s/jobs/start-next/accepted", thingName),
		fmt.Sprintf("$aws/things/%s/jobs/start-next/rejected", thingName),
	}
	pubTopic := fmt.Sprintf("$aws/things/%s/jobs/start-next", thingName)

	payload, err := json.Marshal(req)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, client.Timeouts)
	defer cancel()
	return handleAsync[StartNextPendingJobExecutionOutput](ctx, client.mc, payload, topics, pubTopic)
}

// DescribeJobExecution gets detailed information about a job execution.
func (client *Client) DescribeJobExecution(ctx context.Context, thingName string, jobId string, req DescribeJobExecutionInput) (ret DescribeJobExecutionOutput, err error) {
	topics := []string{
		fmt.Sprintf("$aws/things/%s/jobs/%s/get/accepted", thingName, jobId),
		fmt.Sprintf("$aws/things/%s/jobs/%s/get/rejected", thingName, jobId),
	}
	pubTopic := fmt.Sprintf("$aws/things/%s/jobs/%s/get", thingName, jobId)

	payload, err := json.Marshal(req)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, client.Timeouts)
	defer cancel()
	return handleAsync[DescribeJobExecutionOutput](ctx, client.mc, payload, topics, pubTopic)
}

// UpdateJobExecution updates the status of a job execution.
func (client *Client) UpdateJobExecution(ctx context.Context, thingName string, jobId string, req UpdateJobExecutionInput) (ret iotjobsdataplane.UpdateJobExecutionOutput, err error) {
	topics := []string{
		fmt.Sprintf("$aws/things/%s/jobs/%s/update/accepted", thingName, jobId),
		fmt.Sprintf("$aws/things/%s/jobs/%s/update/rejected", thingName, jobId),
	}
	pubTopic := fmt.Sprintf("$aws/things/%s/jobs/%s/update", thingName, jobId)

	payload, err := json.Marshal(req)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, client.Timeouts)
	defer cancel()
	return handleAsync[iotjobsdataplane.UpdateJobExecutionOutput](ctx, client.mc, payload, topics, pubTopic)
}
