// These contents are copied and slightly modified from aws-sdk-go-v2
// https://github.com/aws/aws-sdk-go-v2/tree/main/service/iotjobsdataplane
// SPDX-License-Identifier: Apache-2.0
package jobs

import "github.com/aws/aws-sdk-go-v2/service/iotjobsdataplane/types"

type UpdateJobExecutionInput struct {
	// The new status for the job execution (IN_PROGRESS, FAILED, SUCCESS, or
	// REJECTED). This must be specified on every update.
	//
	// This member is required.
	Status types.JobExecutionStatus `json:"status"`

	// Optional. A number that identifies a particular job execution on a particular
	// device.
	ExecutionNumber *int64 `json:"executionNumber"`

	// Optional. The expected current version of the job execution. Each time you
	// update the job execution, its version is incremented. If the version of the job
	// execution stored in Jobs does not match, the update is rejected with a
	// VersionMismatch error, and an ErrorResponse that contains the current job
	// execution status data is returned. (This makes it unnecessary to perform a
	// separate DescribeJobExecution request in order to obtain the job execution
	// status data.)
	ExpectedVersion *int64 `json:"expectedVersion"`

	// Optional. When set to true, the response contains the job document. The default
	// is false.
	IncludeJobDocument *bool `json:"includeJobDocument"`

	// Optional. When included and set to true, the response contains the
	// JobExecutionState data. The default is false.
	IncludeJobExecutionState *bool `json:"includeJobExecutionState"`

	// Optional. A collection of name/value pairs that describe the status of the job
	// execution. If not specified, the statusDetails are unchanged.
	StatusDetails map[string]string `json:"statusDetails"`

	// Specifies the amount of time this device has to finish execution of this job. If
	// the job execution status is not set to a terminal state before this timer
	// expires, or before the timer is reset (by again calling UpdateJobExecution,
	// setting the status to IN_PROGRESS and specifying a new timeout value in this
	// field) the job execution status will be automatically set to TIMED_OUT. Note
	// that setting or resetting this timeout has no effect on that job execution
	// timeout which may have been specified when the job was created (CreateJob using
	// field timeoutConfig).
	StepTimeoutInMinutes *int64 `json:"stepTimeoutInMinutes"`

	ClientToken string `json:"clientToken"`
	Timestamp   int    `json:"timestamp"`
}
