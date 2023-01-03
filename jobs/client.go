// SPDX-License-Identifier: Apache-2.0
package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/iotjobsdataplane"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/shirou/aws-iot-device-lib/internal/mqttutils"
)

type Client struct {
	mc mqtt.Client
}

func NewClient(mc mqtt.Client) (*Client, error) {
	client := &Client{
		mc: mc,
	}

	return client, nil
}

type getPendingJobExecutionsRequest struct {
	ClientToken string `json:"clientToken"`
}

// GetPendingJobExecutions gets detailed information about a job execution.
func (client *Client) GetPendingJobExecutions(ctx context.Context, thingName string, req iotjobsdataplane.GetPendingJobExecutionsInput) (ret iotjobsdataplane.GetPendingJobExecutionsOutput, err error) {
	topics := []string{
		fmt.Sprintf("$aws/things/%s/jobs/get/accepted", thingName),
		fmt.Sprintf("$aws/things/%s/jobs/get/rejected", thingName),
	}

	accepted := make(chan iotjobsdataplane.GetPendingJobExecutionsOutput)
	rejected := make(chan error)

	callback := func(mc mqtt.Client, msg mqtt.Message) {
		if err := IsError(msg.Payload()); err != nil {
			rejected <- err
			return
		}
		var output iotjobsdataplane.GetPendingJobExecutionsOutput
		if err = json.Unmarshal(msg.Payload(), &output); err != nil {
			rejected <- err
			return
		}

		if strings.HasSuffix(msg.Topic(), "/jobs/get/accepted") {
			accepted <- output
		} else if strings.HasSuffix(msg.Topic(), "/jobs/get/rejected") {
			rejected <- fmt.Errorf("rejected") // TODO
		} else {
			rejected <- fmt.Errorf("unknown topic subscribed, %s", msg.Topic())
			return
		}
	}

	if err = mqttutils.Subscribe(client.mc, topics, 0, callback); err != nil {
		return
	}
	defer func() {
		close(accepted)
		close(rejected)
		err = mqttutils.JoinErrors(err, mqttutils.Unsubscribe(client.mc, topics))
	}()

	pubTopic := fmt.Sprintf("$aws/things/%s/jobs/get", thingName)
	payload, err := json.Marshal(req)
	if err != nil {
		return
	}
	if err = mqttutils.Publish(client.mc, pubTopic, 0, payload); err != nil {
		return
	}

	for {
		select {
		case r := <-accepted:
			return r, nil
		case <-rejected:
			return ret, fmt.Errorf("rejected")
		case <-ctx.Done():
			err = ctx.Err()
			return
		}
	}
}

// StartNextPendingJobExecution gets and starts the next pending job execution for a thing
func (client *Client) StartNextPendingJobExecution(ctx context.Context, thingName string, req iotjobsdataplane.StartNextPendingJobExecutionInput) (ret StartNextPendingJobExecutionOutput, err error) {
	topics := []string{
		fmt.Sprintf("$aws/things/%s/jobs/start-next/accepted", thingName),
		fmt.Sprintf("$aws/things/%s/jobs/start-next/rejected", thingName),
	}

	accepted := make(chan StartNextPendingJobExecutionOutput)
	rejected := make(chan error)

	callback := func(mc mqtt.Client, msg mqtt.Message) {
		if err := IsError(msg.Payload()); err != nil {
			rejected <- err
			return
		}
		var output StartNextPendingJobExecutionOutput
		if err = json.Unmarshal(msg.Payload(), &output); err != nil {
			rejected <- err
			return
		}

		if strings.HasSuffix(msg.Topic(), "accepted") {
			accepted <- output
		} else if strings.HasSuffix(msg.Topic(), "rejected") {
			rejected <- fmt.Errorf("rejected")
		} else {
			rejected <- fmt.Errorf("unknown topic subscribed, %s", msg.Topic())
			return
		}
	}

	if err = mqttutils.Subscribe(client.mc, topics, 0, callback); err != nil {
		return
	}
	defer func() {
		close(accepted)
		close(rejected)
		err = mqttutils.JoinErrors(err, mqttutils.Unsubscribe(client.mc, topics))
	}()

	pubTopic := fmt.Sprintf("$aws/things/%s/jobs/start-next", thingName)
	payload, err := json.Marshal(req)
	if err != nil {
		return
	}
	if err = mqttutils.Publish(client.mc, pubTopic, 0, payload); err != nil {
		return
	}

	for {
		select {
		case r := <-accepted:
			return r, nil
		case err := <-rejected:
			return ret, fmt.Errorf("rejected, %v", err)
		case <-ctx.Done():
			err = ctx.Err()
			return
		}
	}
}

// DescribeJobExecution gets detailed information about a job execution.
func (client *Client) DescribeJobExecution(ctx context.Context, thingName string, jobId string, req DescribeJobExecutionInput) (ret DescribeJobExecutionOutput, err error) {
	topics := []string{
		fmt.Sprintf("$aws/things/%s/jobs/%s/get/accepted", thingName, jobId),
		fmt.Sprintf("$aws/things/%s/jobs/%s/get/rejected", thingName, jobId),
	}

	accepted := make(chan DescribeJobExecutionOutput)
	rejected := make(chan error)

	callback := func(mc mqtt.Client, msg mqtt.Message) {
		if err := IsError(msg.Payload()); err != nil {
			rejected <- err
			return
		}
		var output DescribeJobExecutionOutput
		if err = json.Unmarshal(msg.Payload(), &output); err != nil {
			rejected <- err
			return
		}

		if strings.HasSuffix(msg.Topic(), "accepted") {
			accepted <- output
		} else if strings.HasSuffix(msg.Topic(), "rejected") {
			rejected <- fmt.Errorf("rejected")
		} else {
			err = fmt.Errorf("unknown topic subscribed, %s", msg.Topic())
			return
		}
	}

	if err = mqttutils.Subscribe(client.mc, topics, 1, callback); err != nil {
		return ret, err
	}
	defer func() {
		close(accepted)
		close(rejected)
		err = mqttutils.JoinErrors(err, mqttutils.Unsubscribe(client.mc, topics))
	}()

	payload, err := json.Marshal(req)
	if err != nil {
		return
	}
	pubTopic := fmt.Sprintf("$aws/things/%s/jobs/%s/get", thingName, jobId)
	if err = mqttutils.Publish(client.mc, pubTopic, 1, payload); err != nil {
		return ret, err
	}

	for {
		select {
		case r := <-accepted:
			return r, nil
		case err = <-rejected:
			return ret, err
		case <-ctx.Done():
			err = ctx.Err()
			return ret, err
		}
	}
}

func (client *Client) UpdateJobExecution(ctx context.Context, thingName string, jobId string, req UpdateJobExecutionInput) (ret iotjobsdataplane.UpdateJobExecutionOutput, err error) {
	topics := []string{
		fmt.Sprintf("$aws/things/%s/jobs/%s/update/accepted", thingName, jobId),
		fmt.Sprintf("$aws/things/%s/jobs/%s/update/rejected", thingName, jobId),
	}

	accepted := make(chan iotjobsdataplane.UpdateJobExecutionOutput)
	rejected := make(chan error)

	callback := func(mc mqtt.Client, msg mqtt.Message) {
		if err := IsError(msg.Payload()); err != nil {
			rejected <- err
			return
		}
		var output iotjobsdataplane.UpdateJobExecutionOutput
		if err = json.Unmarshal(msg.Payload(), &output); err != nil {
			rejected <- err
			return
		}

		if strings.HasSuffix(msg.Topic(), "accepted") {
			accepted <- output
		} else if strings.HasSuffix(msg.Topic(), "rejected") {
			rejected <- fmt.Errorf("rejected")
		} else {
			err = fmt.Errorf("unknown topic subscribed, %s", msg.Topic())
			return
		}
	}

	if err = mqttutils.Subscribe(client.mc, topics, 0, callback); err != nil {
		return
	}
	defer func() {
		close(accepted)
		close(rejected)
		err = mqttutils.JoinErrors(err, mqttutils.Unsubscribe(client.mc, topics))
	}()

	payload, err := json.Marshal(req)
	if err != nil {
		return
	}
	pubTopic := fmt.Sprintf("$aws/things/%s/jobs/%s/update", thingName, jobId)
	if err = mqttutils.Publish(client.mc, pubTopic, 0, payload); err != nil {
		return
	}

	for {
		select {
		case r := <-accepted:
			return r, nil
		case err = <-rejected:
			return
		case <-ctx.Done():
			err = ctx.Err()
			return
		}
	}
}
