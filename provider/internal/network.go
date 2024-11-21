// Copyright 2024, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import "github.com/pulumi/pulumi-go-provider/infer"

var _ = (infer.Enum[NetworkMode])((*NetworkMode)(nil))

// NetworkMode is the --network parameter for a build.
type NetworkMode string

const (
	// Default network mode.
	Default NetworkMode = "default"
	// Host network mode.
	Host NetworkMode = "host"
	// None or no network mode.
	None NetworkMode = "none"
)

// Values returns all valid NetworkMode values for SDK generation.
func (NetworkMode) Values() []infer.EnumValue[NetworkMode] {
	return []infer.EnumValue[NetworkMode]{
		{
			Value:       Default,
			Description: "The default sandbox network mode.",
		},
		{
			Value:       Host,
			Description: "Host network mode.",
		},
		{
			Value:       None,
			Description: "Disable network access.",
		},
	}
}

func (n *NetworkMode) String() string {
	if n == nil {
		return string(Default)
	}
	return string(*n)
}
