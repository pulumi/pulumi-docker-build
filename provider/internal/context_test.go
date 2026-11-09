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

import (
	"bufio"
	gofs "io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _dockerfile = "Dockerfile"

func TestValidateContext(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		c       *BuildContext
		givenD  Dockerfile
		preview bool

		wantD   *Dockerfile
		wantC   *Context
		wantErr string
	}{
		{
			name: "relative",
			c: &BuildContext{Context: Context{
				Location: "../internal/../internal/testdata/noop",
			}},
			wantD: &Dockerfile{
				Location: "../internal/testdata/noop/Dockerfile",
			},
		},
		{
			name: "missing directory",
			c: &BuildContext{Context: Context{
				Location: "/does/not/exist/",
			}},
			wantErr: "not a valid directory",
		},
		{
			name: "missing default Dockerfile",
			c: &BuildContext{Context: Context{
				Location: "testdata",
			}},
			wantD: &Dockerfile{Location: "testdata/Dockerfile"},
		},
		{
			name: "with explicit Dockerfile",
			c: &BuildContext{Context: Context{
				Location: "testdata",
			}},
			givenD: Dockerfile{
				Location: "testdata/Dockerfile.invalid",
			},
		},
		{
			name:  "default location",
			c:     &BuildContext{Context: Context{}},
			wantD: &Dockerfile{Location: "Dockerfile"},
		},
		{
			name: "remote context doesn't default to local Dockerfile",
			c: &BuildContext{Context: Context{
				Location: "https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile",
			}},
			wantD: &Dockerfile{},
		},
		{
			name:    "preview",
			c:       &BuildContext{Context: Context{}},
			preview: true,
		},
		{
			name:  "missing context defaults to current directory",
			c:     nil,
			wantC: &Context{Location: "."},
			wantD: &Dockerfile{Location: "Dockerfile"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d, c, err := tt.c.validate(tt.preview, &tt.givenD)

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErr)
			}

			if tt.wantD != nil {
				assert.Equal(t, tt.wantD.Location, d.Location)
				assert.Equal(t, tt.wantD.Inline, d.Inline)
			}
			if tt.wantC != nil {
				assert.Equal(t, tt.wantC.Location, c.Location)
			}
		})
	}
}

func TestHashIgnoresFile(t *testing.T) {
	t.Parallel()

	step1Dir := "./testdata/ignores/basedir"
	baseResult, err := hashBuildContext(step1Dir, filepath.Join(step1Dir, _dockerfile), nil)
	require.NoError(t, err)

	step2Dir := "./testdata/ignores/basedir-with-ignored-files"
	result, err := hashBuildContext(step2Dir, filepath.Join(step2Dir, _dockerfile), nil)
	require.NoError(t, err)

	assert.Equal(t, result, baseResult)
}

// Tests that we handle .dockerignore exclusions such as "!foo/*/bar".
//
// See:
// - https://github.com/moby/moby/issues/30018
// - https://github.com/moby/moby/issues/45608
//
// Buildkit handles these correctly (according to spec), Docker's classic builder does not.
func TestHashIgnoresWildcards(t *testing.T) {
	t.Parallel()
	baselineDir := "testdata/ignores-wildcard/basedir"
	baselineResult, err := hashBuildContext(
		baselineDir,
		filepath.Join(baselineDir, _dockerfile),
		nil,
	)
	require.NoError(t, err)

	modIgnoredDir := "testdata/ignores-wildcard/basedir-modified-ignored-file"
	modIgnoredResult, err := hashBuildContext(
		modIgnoredDir,
		filepath.Join(modIgnoredDir, _dockerfile),
		nil,
	)
	require.NoError(t, err)

	modIncludedDir := "testdata/ignores-wildcard/basedir-modified-included-file"
	modIncludedResult, err := hashBuildContext(
		modIncludedDir,
		filepath.Join(modIncludedDir, _dockerfile),
		nil,
	)
	require.NoError(t, err)

	assert.Equal(
		t,
		baselineResult,
		modIgnoredResult,
		"hash should not change when modifying ignored files",
	)
	assert.NotEqual(t, baselineResult, modIncludedResult,
		"hash should change when modifying included (via wildcard ignore exclusion) files")
}

func BenchmarkHashBuildContext(b *testing.B) {
	dir := "testdata/ignores-wildcard/basedir-modified-ignored-file"
	for n := 0; n < b.N; n++ {
		_, err := hashBuildContext(dir, filepath.Join(dir, _dockerfile), nil)
		require.NoError(b, err)

	}
}

// Tests that we handle .dockerignore exclusions such as "!foo/*/bar", as above, when using a
// relative context path.
//
//nolint:paralleltest // Incompatible with os.Chdir.
func TestHashIgnoresWildcardsRelative(t *testing.T) {
	err := os.Chdir("testdata")
	require.NoError(t, err)
	defer func() {
		err = os.Chdir("..")
		require.NoError(t, err)
	}()

	baselineDir := "../testdata/ignores-wildcard/basedir"
	baselineResult, err := hashBuildContext(
		baselineDir,
		filepath.Join(baselineDir, _dockerfile),
		nil,
	)
	require.NoError(t, err)

	modIgnoredDir := "../testdata/ignores-wildcard/basedir-modified-ignored-file"
	modIgnoredResult, err := hashBuildContext(
		modIgnoredDir,
		filepath.Join(modIgnoredDir, _dockerfile),
		nil,
	)
	require.NoError(t, err)

	modIncludedDir := "../testdata/ignores-wildcard/basedir-modified-included-file"
	modIncludedResult, err := hashBuildContext(
		modIncludedDir,
		filepath.Join(modIncludedDir, _dockerfile),
		nil,
	)
	require.NoError(t, err)

	assert.Equal(
		t,
		baselineResult,
		modIgnoredResult,
		"hash should not change when modifying ignored files",
	)
	assert.NotEqual(t, baselineResult, modIncludedResult,
		"hash should change when modifying included (via wildcard ignore exclusion) files")
}

func TestHashIgnoresDockerfileOutsideDirMove(t *testing.T) {
	t.Parallel()
	appDir := "./testdata/dockerfile-location-irrelevant/app"
	baseResult, err := hashBuildContext(
		appDir,
		"./testdata/dockerfile-location-irrelevant/step1.Dockerfile",
		nil,
	)
	require.NoError(t, err)

	result, err := hashBuildContext(
		appDir,
		"./testdata/dockerfile-location-irrelevant/step2.Dockerfile",
		nil,
	)
	require.NoError(t, err)

	assert.Equal(t, result, baseResult)
}

func TestHashRenamingMatters(t *testing.T) {
	t.Parallel()
	step1Dir := "./testdata/filemode-matters/step1"
	baseResult, err := hashBuildContext(step1Dir, filepath.Join(step1Dir, _dockerfile), nil)
	require.NoError(t, err)

	step2Dir := "./testdata/renaming-matters/step2"
	result, err := hashBuildContext(step2Dir, filepath.Join(step2Dir, _dockerfile), nil)
	require.NoError(t, err)

	assert.NotEqual(t, result, baseResult)
}

func TestHashFilemodeMatters(t *testing.T) {
	t.Parallel()
	step1Dir := "./testdata/filemode-matters/step1"
	baseResult, err := hashBuildContext(step1Dir, filepath.Join(step1Dir, _dockerfile), nil)
	require.NoError(t, err)

	step2Dir := "./testdata/filemode-matters/step2-chmod-x"
	result, err := hashBuildContext(step2Dir, filepath.Join(step2Dir, _dockerfile), nil)
	require.NoError(t, err)

	assert.NotEqual(t, result, baseResult)
}

func TestHashIgnoresUmaskDerivedPermissions(t *testing.T) {
	t.Parallel()
	// git records only 100644/100755, so a checkout's group/other permission bits come from the
	// umask of whoever materialized the files: umask 022 yields 0644, umask 002 yields 0664. Two
	// checkouts of the same commit must therefore hash identically, or contextHash reports a
	// change that no commit can explain and the diff never converges.
	write := func(t *testing.T, mode gofs.FileMode) string {
		t.Helper()
		dir := t.TempDir()
		dockerfile := filepath.Join(dir, "Dockerfile")
		require.NoError(t, os.WriteFile(dockerfile, []byte("FROM scratch\nCOPY app.txt /\n"), mode))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "app.txt"), []byte("content"), mode))
		// os.WriteFile is subject to the process umask, so set the modes explicitly.
		require.NoError(t, os.Chmod(dockerfile, mode))
		require.NoError(t, os.Chmod(filepath.Join(dir, "app.txt"), mode))
		h, err := hashBuildContext(dir, dockerfile, nil)
		require.NoError(t, err)
		return h
	}

	umask022 := write(t, 0o644)
	umask002 := write(t, 0o664)
	assert.Equal(t, umask022, umask002,
		"group-writable checkout must hash the same as a umask-022 checkout")
}

func TestHashDockerfileErrorIsNotSwallowed(t *testing.T) {
	t.Parallel()
	// A Dockerfile that stats as a local file but cannot be read must surface an error. Returning
	// an empty hash with a nil error makes every subsequent Diff see a changed contextHash.
	dir := t.TempDir()
	dockerfile := filepath.Join(dir, "Dockerfile")
	require.NoError(t, os.WriteFile(dockerfile, []byte("FROM scratch\n"), 0o644))
	require.NoError(t, os.Chmod(dockerfile, 0o000))
	t.Cleanup(func() { _ = os.Chmod(dockerfile, 0o644) })

	if _, err := os.ReadFile(dockerfile); err == nil {
		t.Skip("running as a user that bypasses file permissions")
	}

	_, err := hashBuildContext(dir, dockerfile, nil)
	assert.Error(t, err, "an unreadable Dockerfile must not hash to an empty string with nil error")
}

func TestHashDeepSymlinks(t *testing.T) {
	t.Parallel()
	dir := "./testdata/symlinks"
	_, err := hashBuildContext(dir, filepath.Join(dir, "Dockerfile"), nil)
	assert.NoError(t, err)
}

func TestIgnoreIrregularFiles(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	// Create a Dockerfile
	dockerfile := filepath.Join(dir, "Dockerfile")
	err := os.WriteFile(dockerfile, []byte{}, 0o600)
	require.NoError(t, err)

	// Create a pipe which should be ignored. (We will time out trying to read
	// it if it's not.)
	pipe := filepath.Join(dir, "pipe")
	err = syscall.Mkfifo(pipe, 0o666)
	require.NoError(t, err)
	// Confirm it's irregular.
	fi, err := os.Stat(pipe)
	require.NoError(t, err)
	assert.False(t, fi.Mode().IsRegular())

	_, err = hashBuildContext(dir, dockerfile, nil)
	assert.NoError(t, err)
}

func TestHashUnignoredDirs(t *testing.T) {
	t.Parallel()
	step1Dir := "./testdata/unignores/basedir"
	baseResult, err := hashBuildContext(step1Dir, filepath.Join(step1Dir, _dockerfile), nil)
	require.NoError(t, err)

	step2Dir := "./testdata/unignores/basedir-with-unignored-files"
	unignoreResult, err := hashBuildContext(step2Dir, filepath.Join(step2Dir, _dockerfile), nil)
	require.NoError(t, err)

	assert.Equal(t, baseResult, unignoreResult)
}

func TestDockerIgnore(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string

		dockerfile string
		context    string
		fs         map[string]string

		want    []string
		wantErr error
	}{
		{
			name:       "Dockerfile with root dockerignore",
			dockerfile: fooDockerfilePath,
			fs: map[string]string{
				dockerignoreName: rootIgnore,
			},
			want: []string{rootIgnore},
		},
		{
			name:       "Dockerfile with root dockerignore and custom dockerignore",
			dockerfile: fooDockerfilePath,
			fs: map[string]string{
				"foo/Dockerfile.dockerignore": customIgnore,
				dockerignoreName:              rootIgnore,
			},
			want: []string{customIgnore},
		},
		{
			name:       "Dockerfile with root dockerignore and relative context",
			dockerfile: fooDockerfilePath,
			context:    "../",
			fs: map[string]string{
				"../.dockerignore": rootIgnore,
			},
			want: []string{rootIgnore},
		},
		{
			name:       "Dockerfile without root dockerignore",
			dockerfile: fooDockerfilePath,
			want:       nil,
		},
		{
			name:       "Dockerfile with invalid root dockerignore",
			dockerfile: fooDockerfilePath,
			fs: map[string]string{
				dockerignoreName: strings.Repeat("*", bufio.MaxScanTokenSize),
			},
			wantErr: bufio.ErrTooLong,
		},
		{
			name:       "custom.Dockerfile without custom dockerignore and without root dockerignore",
			dockerfile: "./foo/custom.Dockerfile",
			want:       nil,
		},
		{
			name:       "custom.Dockerfile with custom dockerignore and without root dockerignore",
			dockerfile: "./foo/custom.Dockerfile",
			fs: map[string]string{
				"foo/custom.Dockerfile.dockerignore": customIgnore,
			},
			want: []string{customIgnore},
		},
		{
			name:       "custom.Dockerfile with custom dockerignore and with root dockerignore",
			dockerfile: "foo/custom.Dockerfile",
			fs: map[string]string{
				"foo/custom.Dockerfile.dockerignore": customIgnore,
				dockerignoreName:                     rootIgnore,
			},
			want: []string{customIgnore},
		},
		{
			name:       "custom.Dockerfile without custom dockerignore and with root dockerignore",
			dockerfile: "foo/custom.Dockerfile",
			fs: map[string]string{
				dockerignoreName: rootIgnore,
			},
			want: []string{rootIgnore},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fs := afero.NewMemMapFs()
			for fname, fdata := range tt.fs {
				f, err := fs.Create(fname)
				require.NoError(t, err)
				_, err = f.WriteString(fdata)
				require.NoError(t, err)
			}
			actual, err := getIgnorePatterns(fs, tt.dockerfile, tt.context)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, actual)
		})
	}
}
