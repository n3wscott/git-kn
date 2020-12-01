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

// repo.go provides generic functions related to Repo

package ghutil

import (
	"fmt"

	"github.com/google/go-github/v32/github"
)

// ListRepos lists repos under org
func (gc *GithubClient) ListRepos(orgs ...string) ([]Repo, error) {
	repos := make([]Repo, 0)
	for _, org := range orgs {
		repoListOptions := &github.RepositoryListOptions{}
		genericList, err := gc.depaginate(
			"listing repos",
			maxRetryCount,
			&repoListOptions.ListOptions,
			func() ([]interface{}, *github.Response, error) {
				page, resp, err := gc.Client.Repositories.List(ctx, org, repoListOptions)
				var interfaceList []interface{}
				if nil == err {
					for _, repo := range page {
						interfaceList = append(interfaceList, repo)
					}
				}
				return interfaceList, resp, err
			},
		)
		if err != nil {
			return nil, err
		}

		for _, elem := range genericList {
			r := elem.(*github.Repository)
			if r.Archived != nil && !*r.Archived {
				repo := Repo{
					Org:   org,
					Name:  r.GetName(),
					URL:   r.GetHTMLURL(),
					Stars: r.GetStargazersCount(),
					Forks: r.GetForksCount(),
				}
				//fmt.Println(r)
				repos = append(repos, repo)
			}
		}
	}
	return repos, nil
}

// ListBranches lists branchs for given repo
func (gc *GithubClient) ListBranches(org, repo string) ([]*github.Branch, error) {
	genericList, err := gc.depaginate(
		fmt.Sprintf("listing Pull request from org %q and base %q", org, repo),
		maxRetryCount,
		&github.ListOptions{},
		func() ([]interface{}, *github.Response, error) {
			page, resp, err := gc.Client.Repositories.ListBranches(ctx, org, repo, nil)
			var interfaceList []interface{}
			if nil == err {
				for _, PR := range page {
					interfaceList = append(interfaceList, PR)
				}
			}
			return interfaceList, resp, err
		},
	)
	res := make([]*github.Branch, len(genericList))
	for i, elem := range genericList {
		res[i] = elem.(*github.Branch)
	}
	return res, err
}

type Repo struct {
	Org  string
	Name string
	URL  string

	Stars int
	Forks int

	Fork string
}

// ListRepos lists repos under org
func (gc *GithubClient) JoinRepos(fork string, orgs ...string) ([]Repo, error) {
	forks, err := gc.listRepos(fork)
	if err != nil {
		return nil, err
	}

	// parent to fork
	ptf := map[string]*github.Repository{}
	for _, fork := range forks {
		fmt.Println(fork)

		if fork.Parent != nil {
			fmt.Print("has parent: ", fork.Parent)
		}

		if fork.Fork != nil && *fork.Fork && fork.ForksURL != nil && fork.URL != nil {
			ptf[*fork.ForksURL] = fork
		}
	}

	repos := make([]Repo, 0)
	for _, org := range orgs {
		orgRepos, err := gc.listRepos(org)
		if err != nil {
			return nil, err
		}

		for _, or := range orgRepos {
			repo := Repo{
				Org:  org,
				Name: or.GetName(),
				Fork: "",
			}
			if fork, found := ptf[*or.URL]; found {
				repo.Fork = fork.GetName()
			}
			repos = append(repos, repo)
		}
	}
	return repos, err
}

func (gc *GithubClient) listRepos(org string) ([]*github.Repository, error) {
	repoListOptions := &github.RepositoryListOptions{}
	genericList, err := gc.depaginate(
		"listing repos",
		maxRetryCount,
		&repoListOptions.ListOptions,
		func() ([]interface{}, *github.Response, error) {
			page, resp, err := gc.Client.Repositories.List(ctx, org, repoListOptions)
			var interfaceList []interface{}
			if nil == err {
				for _, repo := range page {
					interfaceList = append(interfaceList, repo)
				}
			}
			return interfaceList, resp, err
		},
	)
	res := make([]*github.Repository, 0, len(genericList))
	for _, elem := range genericList {
		r := elem.(*github.Repository)
		if r.Archived != nil && !*r.Archived {
			repo, _, err := gc.Client.Repositories.Get(ctx, org, r.GetName())
			if err != nil {
				res = append(res, r)
			} else {
				res = append(res, repo)
			}
		}
	}
	return res, err
}
