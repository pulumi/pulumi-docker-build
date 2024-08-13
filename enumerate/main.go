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
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var (
	rootFlag  = flag.String("root", ".", "directory to search")
	shardFlag = flag.Int("shard", 0, "shard index to collect tests for")
	nFlag     = flag.Int("n", 1, "the total number of shards")
)

var re = regexp.MustCompile(`^Test[A-Z_]`)

type testf struct {
	path string
	name string
}

func main() {
	flag.Parse()

	root := *rootFlag
	shard := *shardFlag
	n := *nFlag
	if shard >= n {
		log.Fatal("shard must be less than n")
	}
	if shard < 0 || n < 0 {
		log.Fatal("must be non-negative")
	}

	tests := []testf{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Parse the file to find test functions
		fileSet := token.NewFileSet()
		node, err := parser.ParseFile(fileSet, path, nil, 0)
		if err != nil {
			return err
		}
		for _, decl := range node.Decls {
			f, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			name := f.Name.Name
			if !re.MatchString(name) {
				continue
			}
			tests = append(tests, testf{path: filepath.Dir(path), name: name})
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Shuffle the tests.
	for i := range tests {
		j := rand.Intn(i + 1) //nolint:gosec // Not cryptographic.
		tests[i], tests[j] = tests[j], tests[i]
	}

	// Assign tests to our shard.
	paths := []string{}
	names := []string{}
	for idx, test := range tests {
		if idx%n != shard {
			continue
		}
		paths = append(paths, "./"+test.path)
		names = append(names, test.name)
	}

	// De-dupe.
	slices.Sort(paths)
	slices.Sort(names)
	paths = slices.Compact(paths)
	names = slices.Compact(names)

	fmt.Printf("-run ^(%s)$ %s\n", strings.Join(names, "|"), strings.Join(paths, " "))
}
