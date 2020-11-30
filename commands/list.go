/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/n3wscott/git-kn/pkg/config"
)

func addListCmd(root *cobra.Command) {

	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List the Knative repos.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Config()
			gh, err := cfg.GitHubClient()
			if err != nil {
				return err
			}

			// for all given orgs, list the repos.

			repos, err := gh.JoinRepos(cfg.Fork(), cfg.Orgs()...)
			if err != nil {
				return err
			}
			for _, repo := range repos {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf("%s %s %s", repo.Org, repo.Name, repo.Fork))
			}

			return nil
		},
	}

	root.AddCommand(cmd)
}
