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

package deprecated

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

func TestConfigEncoding(t *testing.T) {
	t.Parallel()

	type testCase struct {
		ty    schema.TypeSpec
		given resource.PropertyValue
		want  resource.PropertyValue
	}

	knownKey := "mykey"

	makeEnc := func(typ schema.TypeSpec) *ConfigEncoding {
		return New(
			schema.ConfigSpec{
				Variables: map[string]schema.PropertySpec{
					knownKey: {
						TypeSpec: typ,
					},
				},
			},
		)
	}

	checkUnmarshal := func(t *testing.T, tc testCase) {
		enc := makeEnc(tc.ty)
		key := resource.PropertyKey(knownKey)

		actual, err := enc.unmarshalPropertyValue(key, tc.given)
		require.NoError(t, err)
		assert.Equal(t, tc.want, actual)
	}

	turnaroundTestCases := []testCase{
		{
			schema.TypeSpec{Type: "boolean"},
			resource.NewPropertyValue(`true`),
			resource.NewBoolProperty(true),
		},
		{
			schema.TypeSpec{Type: "boolean"},
			resource.NewPropertyValue(`false`),
			resource.NewBoolProperty(false),
		},
		{
			schema.TypeSpec{Type: "integer"},
			resource.NewPropertyValue(`0`),
			resource.NewNumberProperty(0),
		},
		{
			schema.TypeSpec{Type: "integer"},
			resource.NewPropertyValue(`42`),
			resource.NewNumberProperty(42),
		},
		{
			schema.TypeSpec{Type: "number"},
			resource.NewPropertyValue(`0`),
			resource.NewNumberProperty(0.0),
		},
		{
			schema.TypeSpec{Type: "number"},
			resource.NewPropertyValue(`42.5`),
			resource.NewNumberProperty(42.5),
		},
		{
			schema.TypeSpec{Type: "string"},
			resource.NewStringProperty(""),
			resource.NewStringProperty(""),
		},
		{
			schema.TypeSpec{Type: "string"},
			resource.NewStringProperty("hello"),
			resource.NewStringProperty("hello"),
		},
		{
			schema.TypeSpec{Type: "array"},
			resource.NewPropertyValue(`[]`),
			resource.NewArrayProperty([]resource.PropertyValue{}),
		},
		{
			schema.TypeSpec{Type: "array"},
			resource.NewPropertyValue(`["hello","there"]`),
			resource.NewArrayProperty([]resource.PropertyValue{
				resource.NewStringProperty("hello"),
				resource.NewStringProperty("there"),
			}),
		},
		{
			schema.TypeSpec{Type: "object"},
			resource.NewPropertyValue(`{}`),
			resource.NewObjectProperty(resource.PropertyMap{}),
		},
		{
			schema.TypeSpec{Type: "object"},
			resource.NewPropertyValue(`{"key":"value"}`),
			resource.NewObjectProperty(resource.PropertyMap{
				"key": resource.NewStringProperty("value"),
			}),
		},
	}

	t.Run("turnaround", func(t *testing.T) {
		for i, tc := range turnaroundTestCases {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				t.Parallel()
				checkUnmarshal(t, tc)
			})
		}
	})

	t.Run("zero_values", func(t *testing.T) {
		// Historically the encoding was able to convert empty strings into type-appropriate zero values.
		cases := []testCase{
			{
				schema.TypeSpec{Type: "boolean"},
				resource.NewPropertyValue(""),
				resource.NewBoolProperty(false),
			},
			{
				schema.TypeSpec{Type: "number"},
				resource.NewPropertyValue(""),
				resource.NewNumberProperty(0.),
			},
			{
				schema.TypeSpec{Type: "integer"},
				resource.NewPropertyValue(""),
				resource.NewNumberProperty(0),
			},
			{
				schema.TypeSpec{Type: "string"},
				resource.NewPropertyValue(""),
				resource.NewStringProperty(""),
			},
			{
				schema.TypeSpec{Type: "object"},
				resource.NewPropertyValue(""),
				resource.NewObjectProperty(make(resource.PropertyMap)),
			},
			{
				schema.TypeSpec{Type: "array"},
				resource.NewPropertyValue(""),
				resource.NewArrayProperty([]resource.PropertyValue{}),
			},
		}
		for _, tc := range cases {
			t.Run(fmt.Sprintf("%v", tc.ty), func(t *testing.T) {
				t.Parallel()
				checkUnmarshal(t, tc)
			})
		}
	})

	t.Run("computed", func(t *testing.T) {
		unk := resource.MakeComputed(resource.NewStringProperty(""))

		for i, tc := range turnaroundTestCases {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				t.Parallel()
				// Unknown sentinel unmarshals to a Computed with a type-appropriate zero value.
				checkUnmarshal(t, testCase{
					ty:    tc.ty,
					given: unk,
					want:  resource.MakeComputed(makeEnc(tc.ty).zeroValue(tc.ty.Type)),
				})
			})
		}
	})

	t.Run("secret", func(t *testing.T) {
		// Unmarshalling happens with KeepSecrets=false, replacing them with the underlying values. This case
		// does not need to be tested.
		//
		// Marhalling however supports sending secrets back to the engine, intending to mark values as secret
		// that happen on paths that are declared as secret in the schema. Due to the limitation of the
		// JSON-in-proto-encoding, secrets are communicated imprecisely as an approximation: if any nested
		// element of a property is secret, the entire property is marshalled as secret.

		var secretCases []testCase

		for _, tc := range turnaroundTestCases {
			secretCases = append(secretCases, testCase{
				ty:    tc.ty,
				given: resource.MakeSecret(tc.given),
				want:  resource.MakeSecret(tc.want),
			})
		}

		for i, tc := range secretCases {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				t.Parallel()
				checkUnmarshal(t, tc)
			})
		}

		t.Run("nested secrets", func(t *testing.T) {
			checkUnmarshal(t, testCase{
				schema.TypeSpec{Type: "object"},
				resource.MakeSecret(resource.NewPropertyValue(`{"key":"val"}`)),
				resource.MakeSecret(resource.NewObjectProperty(resource.PropertyMap{
					"key": resource.NewStringProperty("val"),
				})),
			})
		})
	})

	regressUnmarshalTestCases := []testCase{
		{
			schema.TypeSpec{Type: "array"},
			resource.NewPropertyValue(`
			[
			  {
			    "address": "somewhere.org",
			    "password": {
			      "4dabf18193072939515e22adb298388d": "1b47061264138c4ac30d75fd1eb44270",
			      "value": "some-password"
			    },
			    "username": "some-user"
			  }
			]`),
			resource.NewArrayProperty([]resource.PropertyValue{
				resource.NewObjectProperty(resource.PropertyMap{
					"address":  resource.NewStringProperty("somewhere.org"),
					"password": resource.MakeSecret(resource.NewStringProperty("some-password")),
					"username": resource.NewStringProperty("some-user"),
				}),
			}),
		},
	}

	t.Run("regress-unmarshal", func(t *testing.T) {
		for i, tc := range regressUnmarshalTestCases {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				t.Parallel()
				checkUnmarshal(t, tc)
			})
		}
	})
}
