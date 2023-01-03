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
func (client *Client) GetPendingJobExecutions(ctx context.Context, req iotjobsdataplane.GetPendingJobExecutionsInput) (ret iotjobsdataplane.GetPendingJobExecutionsOutput, err error) {
	topics := []string{
		fmt.Sprintf("$aws/things/%s/jobs/get/accepted", *req.ThingName),
		fmt.Sprintf("$aws/things/%s/jobs/get/rejected", *req.ThingName),
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

	pubTopic := fmt.Sprintf("$aws/things/%s/jobs/get", *req.ThingName)
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
func (client *Client) StartNextPendingJobExecution(ctx context.Context, req iotjobsdataplane.StartNextPendingJobExecutionInput) (ret StartNextPendingJobExecutionOutput, err error) {
	topics := []string{
		fmt.Sprintf("$aws/things/%s/jobs/start-next/accepted", *req.ThingName),
		fmt.Sprintf("$aws/things/%s/jobs/start-next/rejected", *req.ThingName),
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

	pubTopic := fmt.Sprintf("$aws/things/%s/jobs/start-next", *req.ThingName)
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
func (client *Client) DescribeJobExecution(ctx context.Context, req DescribeJobExecutionInput) (ret DescribeJobExecutionOutput, err error) {
	topics := []string{
		fmt.Sprintf("$aws/things/%s/jobs/%s/get/accepted", *req.ThingName, *req.JobId),
		fmt.Sprintf("$aws/things/%s/jobs/%s/get/rejected", *req.ThingName, *req.JobId),
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
	pubTopic := fmt.Sprintf("$aws/things/%s/jobs/%s/get", *req.ThingName, *req.JobId)
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

func (client *Client) UpdateJobExecution(ctx context.Context, req UpdateJobExecutionInput) (ret iotjobsdataplane.UpdateJobExecutionOutput, err error) {
	topics := []string{
		fmt.Sprintf("$aws/things/%s/jobs/%s/update/accepted", *req.ThingName, *req.JobId),
		fmt.Sprintf("$aws/things/%s/jobs/%s/update/rejected", *req.ThingName, *req.JobId),
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
	pubTopic := fmt.Sprintf("$aws/things/%s/jobs/%s/update", *req.ThingName, *req.JobId)
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

type JobExecutionsChangedHandler func()

// JobExecutionsChanged sent whenever a job execution is added to or removed from the list of pending job executions for a thing.
func (client *Client) JobExecutionsChanged(ctx context.Context, thingName string, handler JobExecutionsChangedHandler) (err error) {
	topics := []string{fmt.Sprintf("$aws/things/%s/jobs/notify", thingName)}

	callback := func(mc mqtt.Client, msg mqtt.Message) {
		fmt.Println(msg.Topic())
		fmt.Println(string(msg.Payload()))
		handler()
	}
	if err = mqttutils.Subscribe(client.mc, topics, 0, callback); err != nil {
		return
	}
	defer func() {
		err = mqttutils.Unsubscribe(client.mc, topics)
	}()

	return
}

type NextJobExecutionChangedHandler func()

// NextJobExecutionChanged sent whenever there is a change to which job execution is next on the list of pending job executions for a thing
func (client *Client) NextJobExecutionChanged(ctx context.Context, handler NextJobExecutionChangedHandler) (err error) {

	return
}
