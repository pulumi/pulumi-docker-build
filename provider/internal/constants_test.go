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

// Shared string literals used across the package's tests. Hoisted to
// constants to satisfy goconst and keep fixtures consistent.
const (
	fooName = "foo"
	barName = "bar"

	dockerIO          = "docker.io"
	dockerIORef       = "docker.io/foo/bar"
	dockerIOTaggedRef = "docker.io/foo/bar:baz"

	fromScratch = "FROM scratch"
	inlineName  = "inline"
	fooDest     = "/foo"

	contextKey  = "context"
	locationKey = "location"
	exportsKey  = "exports"
	addressKey  = "address"
	knownKey    = "known"

	notReal = "not-real"
	oldBar  = "old_bar"
	newBar  = "new_bar"

	fooDockerfilePath = "./foo/Dockerfile"
	dockerignoreName  = ".dockerignore"
	customIgnore      = "customignore"
	exampleAppContext = "../../examples/app"

	miscName      = "misc"
	awsECRAddress = "1234.dkr.ecr.us-west-2.amazonaws.com"

	rawKey      = "raw"
	usernameKey = "username"
	passwordKey = "password" //nolint:gosec // G101: property-name key, not a real credential.

	rootIgnore             = "rootignore"
	unknownInstructionRUNN = "unknown instruction: RUNN"
	typeTar                = "type=tar"
	testdataNoop           = "testdata/noop"

	// defaultSSHID is buildkit's default SSH agent forward ID (unrelated to
	// NetworkMode.Default).
	defaultSSHID = "default"

	testCacheURL     = "test-cache-url"
	testRuntimeToken = "test-runtime-token" //nolint:gosec // G101: test fixture, not a real credential.
	testResultsURL   = "test-results-url"
)
