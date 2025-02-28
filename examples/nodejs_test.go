//go:build nodejs || all
// +build nodejs all

package examples

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/pulumi/providertest"
	"github.com/pulumi/providertest/optproviderupgrade"
	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/assertpreview"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeExample(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	test := integration.ProgramTestOptions{
		Dir:          path.Join(cwd, "nodejs"),
		Dependencies: []string{"@pulumi/docker-build"},
		Secrets: map[string]string{
			"dockerHubPassword": os.Getenv("DOCKER_HUB_PASSWORD"),
		},
	}

	integration.ProgramTest(t, &test)
}

func TestNodeExampleUpgrade(t *testing.T) {
	t.Parallel()
	var (
		providerName    string = "docker-build"
		baselineVersion string = "0.0.7"
	)

	cwd, err := os.Getwd()
	require.NoError(t, err)

	options := []opttest.Option{
		opttest.DownloadProviderVersion(providerName, baselineVersion),
		opttest.LocalProviderPath(providerName, filepath.Join(cwd, "..", "bin")),
		opttest.YarnLink("@pulumi/docker-build"),
		opttest.TestInPlace(),
	}

	test := pulumitest.NewPulumiTest(t, filepath.Join(cwd, "upgrade-node"), options...)
	result := providertest.PreviewProviderUpgrade(t, test, providerName, baselineVersion,
		optproviderupgrade.DisableAttach())

	assertpreview.HasNoReplacements(t, result)
}

// TestCaching simulates a slow build with --cache-to enabled. We aren't able
// to directly detect cache hits, so we re-run the update and confirm it took
// less time than the image originally took to build.
//
// This is a moderately slow test because we need to "build" (i.e., sleep)
// longer than it would take for cache layer uploads under slow network
// conditions.
func TestCaching(t *testing.T) {
	t.Parallel()

	sleep := 20.0 // seconds

	// Provision ECR outside of our stack, because the cache needs to be shared
	// across updates.
	ecr, ecrOK := tmpEcr(t)

	cwd, err := os.Getwd()
	require.NoError(t, err)

	localCache := t.TempDir()

	tests := []struct {
		name string
		skip bool

		cacheTo   string
		cacheFrom string
		address   string
		username  string
		password  string
	}{
		{
			name:      "local",
			cacheTo:   fmt.Sprintf("type=local,mode=max,oci-mediatypes=true,dest=%s", localCache),
			cacheFrom: fmt.Sprintf("type=local,src=%s", localCache),
		},
		{
			name:      "gha",
			skip:      os.Getenv("ACTIONS_CACHE_URL") == "",
			cacheTo:   "type=gha,mode=max,scope=cache-test",
			cacheFrom: "type=gha,scope=cache-test",
		},
		{
			name:      "dockerhub",
			skip:      os.Getenv("DOCKER_HUB_PASSWORD") == "",
			cacheTo:   "type=registry,mode=max,ref=docker.io/pulumibot/myapp:cache",
			cacheFrom: "type=registry,ref=docker.io/pulumibot/myapp:cache",
			address:   "docker.io",
			username:  "pulumibot",
			password:  os.Getenv("DOCKER_HUB_PASSWORD"),
		},
		{
			name:      "ecr",
			skip:      !ecrOK,
			cacheTo:   fmt.Sprintf("type=registry,mode=max,image-manifest=true,oci-mediatypes=true,ref=%s:cache", ecr.address),
			cacheFrom: fmt.Sprintf("type=registry,ref=%s:cache", ecr.address),
			address:   ecr.address,
			username:  ecr.username,
			password:  ecr.password,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("Missing environment variables")
			}

			sleepFuzzed := sleep + rand.Float64() // Add some fuzz to bust any prior build caches.

			test := integration.ProgramTestOptions{
				Dir: path.Join(cwd, "tests", "caching"),
				ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
					duration, ok := stack.Outputs["durationSeconds"]
					assert.True(t, ok)
					assert.Greater(t, duration.(float64), sleepFuzzed)
				},
				Dependencies: []string{"@pulumi/docker-build"},
				Config: map[string]string{
					"SLEEP_SECONDS": fmt.Sprint(sleepFuzzed),
					"cacheTo":       tt.cacheTo,
					"cacheFrom":     tt.cacheFrom,
					"name":          tt.name,
					"address":       tt.address,
					"username":      tt.username,
				},
				Secrets: map[string]string{
					"password": tt.password,
				},
				NoParallel:  true,
				Quick:       true,
				SkipPreview: true,
				SkipRefresh: true,
				Verbose:     true,
			}

			// First run should be un-cached.
			integration.ProgramTest(t, &test)

			// Now run again and confirm our build was faster due to a cache hit.
			test.ExtraRuntimeValidation = func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				duration, ok := stack.Outputs["durationSeconds"]
				assert.True(t, ok)
				assert.Less(t, duration.(float64), sleepFuzzed)
			}
			test.Config["name"] += "-cached"
			integration.ProgramTest(t, &test)
		})
	}
}

func TestConfig(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	test := integration.ProgramTestOptions{
		Dir:          path.Join(cwd, "tests", "config"),
		Dependencies: []string{"@pulumi/docker-build"},
	}

	integration.ProgramTest(t, &test)
}

type ECR struct {
	address  string
	username string
	password string
}

// tmpEcr creates a new ECR repo and cleans it up after the test concludes.
func tmpEcr(t *testing.T) (ECR, bool) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		return ECR{}, false
	}

	svc := ecr.New(sess)
	name := strings.ToLower(t.Name()) + fmt.Sprint(rand.Intn(1000))

	// Always attempt to delete pre-existing repos, in case our cleanup didn't
	// run.
	_, _ = svc.DeleteRepository(&ecr.DeleteRepositoryInput{
		Force:          aws.Bool(true),
		RepositoryName: aws.String(name),
	})

	params := &ecr.CreateRepositoryInput{
		RepositoryName: aws.String(name),
	}
	resp, err := svc.CreateRepository(params)
	if err != nil {
		return ECR{}, false
	}
	repo := resp.Repository
	t.Cleanup(func() {
		svc.DeleteRepository(&ecr.DeleteRepositoryInput{
			Force:          aws.Bool(true),
			RegistryId:     repo.RegistryId,
			RepositoryName: repo.RepositoryName,
		})
	})

	// Now grab auth for the repo.
	auth, err := svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return ECR{}, false
	}
	b64token := auth.AuthorizationData[0].AuthorizationToken
	token, err := base64.StdEncoding.DecodeString(*b64token)
	if err != nil {
		return ECR{}, false
	}
	parts := strings.SplitN(string(token), ":", 2)
	if len(parts) != 2 {
		return ECR{}, false
	}
	username := parts[0]
	password := parts[1]

	return ECR{
		address:  *repo.RepositoryUri,
		username: username,
		password: password,
	}, true
}
