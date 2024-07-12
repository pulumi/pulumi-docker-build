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

//go:generate go run generate.go yaml ../provider/internal/embed

// Package main ingests a multi-document YAML file and converts it into
// Markdown examples.
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
)

func main() {
	if len(os.Args) < 3 {
		_, _ = fmt.Fprintf(os.Stdout, "Usage: %s <yaml source dir path> <markdown destination path>\n", os.Args[0])
		os.Exit(1)
	}
	yamlPath := os.Args[1]
	mdPath := os.Args[2]

	if !filepath.IsAbs(yamlPath) {
		cwd, err := os.Getwd()
		contract.AssertNoErrorf(err, "getting working directory")
		yamlPath = filepath.Join(cwd, yamlPath)
	}

	if err := os.MkdirAll(mdPath, 0o750); err != nil {
		panic(err)
	}
	fileInfo, err := os.Lstat(mdPath)
	if err != nil || !fileInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "Expect markdown destination %q to be a directory\n", mdPath)
		os.Exit(1)
	}

	yamlFiles, err := os.ReadDir(yamlPath)
	if err != nil {
		panic(err)
	}
	for _, yamlFile := range yamlFiles {
		if err := processYaml(filepath.Join(yamlPath, yamlFile.Name()), mdPath); err != nil {
			log.Fatal(fmt.Errorf("processing %q: %w", yamlFile.Name(), err))
		}
	}
}

func markdownExamples(examples []string) string {
	s := "{{% examples %}}\n## Example Usage\n"
	for _, example := range examples {
		s += example
	}
	s += "{{% /examples %}}"
	return s
}

func markdownExample(description string,
	typescript string,
	python string,
	csharp string,
	golang string,
	yaml string,
	java string,
) string {
	return fmt.Sprintf("{{%% example %%}}\n### %s\n\n"+
		"```typescript\n%s```\n"+
		"```python\n%s```\n"+
		"```csharp\n%s```\n"+
		"```go\n%s```\n"+
		"```yaml\n%s```\n"+
		"```java\n%s```\n"+
		"{{%% /example %%}}\n",
		description, typescript, python, csharp, golang, yaml, java)
}

func convert(language, tempDir, programFile string) (string, error) {
	exampleDir := filepath.Join(tempDir, "example"+language)
	//nolint:gosec // No user-provided input.
	cmd := exec.Command(
		"pulumi",
		"convert",
		"--language",
		language,
		"--out",
		filepath.Clean(filepath.Join(tempDir, exampleDir)),
		"--generate-only",
	)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("converting: %w", err)
	}
	content, err := os.ReadFile(filepath.Clean(filepath.Join(tempDir, exampleDir, programFile)))
	if err != nil {
		return "", fmt.Errorf("reading: %w", err)
	}

	return string(content), nil
}

func processYaml(path, mdDir string) error {
	yamlFile, err := os.Open(filepath.Clean(path))
	if err != nil {
		return err
	}

	base := filepath.Base(path)
	md := strings.NewReplacer(".yaml", ".md", ".yml", ".md").Replace(base)

	defer contract.IgnoreClose(yamlFile)
	decoder := yaml.NewDecoder(yamlFile)
	exampleStrings := []string{}
	for {
		keepGoing, err := func() (bool, error) {
			example := map[string]interface{}{}
			err := decoder.Decode(&example)
			if err == io.EOF {
				return false, nil
			}

			description, ok := example["description"].(string)
			if !ok {
				description = ""
			}
			dir, err := os.MkdirTemp("", "")
			if err != nil {
				return false, err
			}

			defer func() {
				contract.IgnoreError(os.RemoveAll(dir))
			}()

			src, err := os.OpenFile(filepath.Clean(filepath.Join(dir, "Pulumi.yaml")), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
			if err != nil {
				return false, err
			}

			fmt.Println("Converting:", example)

			if err := yaml.NewEncoder(src).Encode(example); err != nil {
				return false, err
			}
			contract.AssertNoErrorf(src.Close(), "closing")

			typescript, err := convert("typescript", dir, "index.ts")
			if err != nil {
				return false, err
			}
			python, err := convert("python", dir, "__main__.py")
			if err != nil {
				return false, err
			}
			csharp, err := convert("csharp", dir, "Program.cs")
			if err != nil {
				return false, err
			}
			golang, err := convert("go", dir, "main.go")
			if err != nil {
				return false, err
			}
			java, err := convert("java", dir, "src/main/java/generated_program/App.java")
			if err != nil {
				return false, err
			}

			yamlContent, err := os.ReadFile(filepath.Clean(filepath.Join(dir, "Pulumi.yaml")))
			if err != nil {
				return false, err
			}
			yaml := string(yamlContent)

			exampleStrings = append(exampleStrings, markdownExample(description, typescript, python, csharp, golang, yaml, java))

			return true, nil
		}()
		if err != nil {
			return err
		}
		if !keepGoing {
			break
		}
	}
	_, _ = fmt.Fprintf(os.Stdout, "Writing %s\n", filepath.Join(mdDir, md))
	f, err := os.OpenFile(filepath.Clean(filepath.Join(mdDir, md)), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(f)
	_, err = f.WriteString(markdownExamples(exampleStrings))
	contract.AssertNoErrorf(err, "writing examples")
	return nil
}
