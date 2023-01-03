// These contents are copied and slightly modified from aws-sdk-go-v2
// https://github.com/aws/aws-sdk-go-v2/tree/main/service/iotjobsdataplane
// SPDX-License-Identifier: Apache-2.0
package jobs

import "github.com/aws/smithy-go/middleware"

type StartNextPendingJobExecutionOutput struct {

	// A JobExecution object.
	Execution *JobExecution `json:"execution"`

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata `json:"-"`

	ClientToken string `json:"clientToken"`
	Timestamp   int    `json:"timestamp"`
}
