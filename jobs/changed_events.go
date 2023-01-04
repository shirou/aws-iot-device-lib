// SPDX-License-Identifier: Apache-2.0
package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/shirou/aws-iot-device-lib/internal/mqttutils"
)

type changedHandlerType[V changedMessageType] interface {
	~func(cli *Client, msg V) error
}
type changedMessageType interface {
	JobExecutionsChangedMessage | NextJobExecutionChangedMessage
}

func handleChanged[K changedHandlerType[V], V changedMessageType](ctx context.Context, client *Client, topics []string, handler K) {
	callback := func(mc mqtt.Client, msg mqtt.Message) {
		var je V
		if err := json.Unmarshal(msg.Payload(), &je); err != nil {
			// TODO: how to log the error?
			return
		}
		go handler(client, je)
	}
	if err := mqttutils.Subscribe(client.mc, topics, 0, callback); err != nil {
		return
	}
	defer func() {
		mqttutils.Unsubscribe(client.mc, topics)
	}()

	// forever
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

type JobExecutionsChangedHandler func(cli *Client, msg JobExecutionsChangedMessage) error

// JobExecutionsChanged sent whenever a job execution is added to or removed from the list of pending job executions for a thing.
func (client *Client) JobExecutionsChanged(ctx context.Context, thingName string, handler JobExecutionsChangedHandler) {
	topics := []string{fmt.Sprintf("$aws/things/%s/jobs/notify", thingName)}

	handleChanged(ctx, client, topics, handler)
}

type NextJobExecutionChangedHandler func(cli *Client, msg NextJobExecutionChangedMessage) error

// NextJobExecutionChanged sent whenever there is a change to which job execution is next on the list of pending job executions for a thing
func (client *Client) NextJobExecutionChanged(ctx context.Context, thingName string, handler NextJobExecutionChangedHandler) {
	topics := []string{fmt.Sprintf("$aws/things/%s/jobs/notify-next", thingName)}

	handleChanged(ctx, client, topics, handler)
}
