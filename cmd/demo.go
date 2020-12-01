package main

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/n3wscott/git-kn/pkg/config"
)

type pepper struct {
	Name     string
	HeatUnit int
	Peppers  int
}

func main() {
	cfg := config.Config()
	gh, err := cfg.GitHubClient()
	if err != nil {
		panic(err)
	}

	// for all given orgs, list the repos.

	repos, err := gh.ListRepos(cfg.Orgs()...)
	if err != nil {
		panic(err)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\u279C {{ .Org | faint }} {{ .Name | cyan }}",
		Inactive: "  {{ .Org | faint }} {{ .Name | cyan }}",
		Selected: "\u2735 {{ .Org | faint }} {{ .Name | red | cyan }}",
		Details: `
--------- Repository ----------
{{ "Org:" | faint }}	{{ .Org }}
{{ "Repo:" | faint }}	{{ .Name }}
{{ "URL:" | faint }}	{{ .URL }}
{{ "Stars:" | faint }}	{{ .Stars }}
{{ "Forks:" | faint }}	{{ .Forks }}
`,
	}

	searcher := func(input string, index int) bool {
		repo := repos[index]
		name := strings.Replace(strings.ToLower(repo.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Repositories",
		Items:     repos,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose number %d: %s\n", i+1, repos[i].Name)
}
