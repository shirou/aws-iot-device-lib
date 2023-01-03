// These contents are copied and slightly modified from aws-sdk-go-v2
// https://github.com/aws/aws-sdk-go-v2/tree/main/service/iotjobsdataplane
// SPDX-License-Identifier: Apache-2.0
package jobs

import (
	"github.com/aws/smithy-go/middleware"
)

type JobExecutionStatus string

// Enum values for JobExecutionStatus
const (
	JobExecutionStatusQueued     JobExecutionStatus = "QUEUED"
	JobExecutionStatusInProgress JobExecutionStatus = "IN_PROGRESS"
	JobExecutionStatusSucceeded  JobExecutionStatus = "SUCCEEDED"
	JobExecutionStatusFailed     JobExecutionStatus = "FAILED"
	JobExecutionStatusTimedOut   JobExecutionStatus = "TIMED_OUT"
	JobExecutionStatusRejected   JobExecutionStatus = "REJECTED"
	JobExecutionStatusRemoved    JobExecutionStatus = "REMOVED"
	JobExecutionStatusCanceled   JobExecutionStatus = "CANCELED"
)

// Values returns all known values for JobExecutionStatus. Note that this can be
// expanded in the future, and so it is only as up to date as the client. The
// ordering of this slice is not guaranteed to be stable across updates.
func (JobExecutionStatus) Values() []JobExecutionStatus {
	return []JobExecutionStatus{
		"QUEUED",
		"IN_PROGRESS",
		"SUCCEEDED",
		"FAILED",
		"TIMED_OUT",
		"REJECTED",
		"REMOVED",
		"CANCELED",
	}
}

type DescribeJobExecutionInput struct {

	// The unique identifier assigned to this job when it was created.
	//
	// This member is required.
	JobId *string `json:"-"`

	// The thing name associated with the device the job execution is running on.
	//
	// This member is required.
	ThingName *string `json:"-"`

	// Optional. A number that identifies a particular job execution on a particular
	// device. If not specified, the latest job execution is returned.
	ExecutionNumber *int64 `json:"executionNumber"`

	// Optional. When set to true, the response contains the job document. The default
	// is false.
	IncludeJobDocument *bool `json:"includeJobDocument"`

	ClientToken string `json:"clientToken"`
	Timestamp   int    `json:"timestamp"`
}

type DescribeJobExecutionOutput struct {
	// Contains data about a job execution.
	Execution *JobExecution `json:"execution"`

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	ClientToken string `json:"clientToken"`
	Timestamp   int    `json:"timestamp"`
}

// Contains data about a job execution.
// https://docs.aws.amazon.com/iot/latest/developerguide/jobs-mqtt-https-api.html#jobs-mqtt-job-execution-data
type JobExecution struct {

	// The estimated number of seconds that remain before the job execution status will
	// be changed to TIMED_OUT.
	ApproximateSecondsBeforeTimedOut *int64

	// A number that identifies a particular job execution on a particular device. It
	// can be used later in commands that return or update job execution information.
	ExecutionNumber *int64 `json:"executionNumber"`

	// The content of the job document.
	JobDocument JobDocument `json:"jobDocument"`

	// The unique identifier you assigned to this job when it was created.
	JobId *string `json:"jobId"`

	// The time, in milliseconds since the epoch, when the job execution was last
	// updated.
	LastUpdatedAt int64 `json:"lastUpdatedAt"`

	// The time, in milliseconds since the epoch, when the job execution was enqueued.
	QueuedAt int64 `json:"queuedAt"`

	// The time, in milliseconds since the epoch, when the job execution was started.
	StartedAt *int64 `json:"startedAt"`

	// The status of the job execution. Can be one of: "QUEUED", "IN_PROGRESS",
	// "FAILED", "SUCCESS", "CANCELED", "REJECTED", or "REMOVED".
	Status JobExecutionStatus `json:"status"`

	// A collection of name/value pairs that describe the status of the job execution.
	StatusDetails map[string]string `json:"statusDetails"`

	// The name of the thing that is executing the job.
	ThingName *string `json:"thingName"`

	// The version of the job execution. Job execution versions are incremented each
	// time they are updated by a device.
	VersionNumber int64 `json:"versionNumber"`

	ClientToken string `json:"clientToken"`
	Timestamp   int    `json:"timestamp"`
}

// Contains data about the state of a job execution.
type JobExecutionState struct {

	// The status of the job execution. Can be one of: "QUEUED", "IN_PROGRESS",
	// "FAILED", "SUCCESS", "CANCELED", "REJECTED", or "REMOVED".
	Status JobExecutionStatus

	// A collection of name/value pairs that describe the status of the job execution.
	StatusDetails map[string]string

	// The version of the job execution. Job execution versions are incremented each
	// time they are updated by a device.
	VersionNumber int64
}

// Contains a subset of information about a job execution.
type JobExecutionSummary struct {

	// A number that identifies a particular job execution on a particular device.
	ExecutionNumber *int64 `json:"executionNumber"`

	// The unique identifier you assigned to this job when it was created.
	JobId *string `json:"jobId"`

	// The time, in milliseconds since the epoch, when the job execution was last
	// updated.
	LastUpdatedAt int64 `json:"lastUpdatedAt"`

	// The time, in milliseconds since the epoch, when the job execution was enqueued.
	QueuedAt int64 `json:"queuedAt"`

	// The time, in milliseconds since the epoch, when the job execution started.
	StartedAt *int64 `json:"startedAt"`

	// The version of the job execution. Job execution versions are incremented each
	// time AWS IoT Jobs receives an update from a device.
	VersionNumber int64 `json:"versionNumber"`
}
