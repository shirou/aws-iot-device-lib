// SPDX-License-Identifier: Apache-2.0
package jobs

import (
	"encoding/json"
	"fmt"
)

// ErrorMessage represents messages if request failed
type ErrorMessage struct {
	ClientToken string `json:"clientToken"`
	Timestamp   int    `json:"timestamp"`
	Code        string `json:"code"`
	Message     string `json:"message"`
}

func IsError(payload []byte) error {
	var msg ErrorMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return nil // This is not a error message format
	}
	if msg.Code == "" {
		return nil
	}

	return fmt.Errorf(msg.Message)
}

type JobExecutions []JobExecution
type JobExecutionsChangedMessage struct {
	Timestamp int                                  `json:"timestamp"`
	Jobs      map[JobExecutionStatus]JobExecutions `json:"jobs"`
}

type NextJobExecutionChangedMessage struct {
	Timestamp int          `json:"timestamp"`
	Execution JobExecution `json:"execution"`
}
