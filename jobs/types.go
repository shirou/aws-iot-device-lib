// SPDX-License-Identifier: Apache-2.0
package jobs

// JobDocument represents JobDocument based on this document.
// https://github.com/awslabs/aws-iot-device-client/tree/main/sample-job-docs
type JobDocument struct {
	Comment string `json:"_comment"`
	Version string `json:"version"`
	Steps   []struct {
		Action struct {
			Name  string `json:"name"`
			Type  string `json:"type"`
			Input struct {
				Handler string   `json:"handler"`
				Args    []string `json:"args"`
				Path    string   `json:"path"`
			} `json:"input"`
			RunAsUser string `json:"runAsUser"`
		} `json:"action"`
	} `json:"steps"`
	FinalStep struct {
		Action struct {
			Name  string `json:"name"`
			Type  string `json:"type"`
			Input struct {
				Handler string   `json:"handler"`
				Args    []string `json:"args"`
				Path    string   `json:"path"`
			} `json:"input"`
			RunAsUser string `json:"runAsUser"`
		} `json:"action"`
	} `json:"finalStep"`
}
