// Copyright 2022 Chainguard, Inc.
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

package cli

import (
	"context"
	"fmt"
	"os"

	"chainguard.dev/melange/pkg/build"
	"github.com/spf13/cobra"
)

func Build() *cobra.Command {
	var buildDate string
	var workspaceDir string
	var pipelineDir string
	var signingKey string
	var useProot bool
	var outDir string

	cmd := &cobra.Command{
		Use:     "build",
		Short:   "Build a package from a YAML configuration file",
		Long:    `Build a package from a YAML configuration file.`,
		Example: `  melange build [config.yaml]`,
		Args:    cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			options := []build.Option{
				build.WithBuildDate(buildDate),
				build.WithWorkspaceDir(workspaceDir),
				build.WithPipelineDir(pipelineDir),
				build.WithSigningKey(signingKey),
				build.WithUseProot(useProot),
				build.WithOutDir(outDir),
			}

			if len(args) > 0 {
				options = append(options, build.WithConfig(args[0]))
			}

			return BuildCmd(cmd.Context(), options...)
		},
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	cmd.Flags().StringVar(&buildDate, "build-date", "", "date used for the timestamps of the files inside the image")
	cmd.Flags().StringVar(&workspaceDir, "workspace-dir", cwd, "directory used for the workspace at /home/build")
	cmd.Flags().StringVar(&pipelineDir, "pipeline-dir", "/usr/share/melange/pipelines", "directory used to store defined pipelines")
	cmd.Flags().StringVar(&signingKey, "signing-key", "", "key to use for signing")
	cmd.Flags().BoolVar(&useProot, "use-proot", false, "whether to use proot for fakeroot")
	cmd.Flags().StringVar(&outDir, "out-dir", cwd, "directory where packages will be output")

	return cmd
}

func BuildCmd(ctx context.Context, opts ...build.Option) error {
	bc, err := build.New(opts...)
	if err != nil {
		return err
	}

	if err := bc.BuildPackage(); err != nil {
		return fmt.Errorf("failed to build package: %w", err)
	}

	return nil
}
