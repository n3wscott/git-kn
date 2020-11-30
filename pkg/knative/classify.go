package knative

import (
	"io"
)

type Project struct {
	Org   string
	Repo  string
	Class string
	Fork  string // name of your fork.

}

var headers = []byte("NAME\tKIND\n")

func ClassifyRepos(repos []string, output io.Writer) ([]Project, error) {
	return nil, nil
	//cfg := config.Config()
	//gh, err := cfg.GitHubClient()
	//forks, err := gh.ListRepos(cfg.Fork())
	//
	//w := printers.GetNewTabWriter(output) // DEBUG
	//defer w.Flush()
	//_, err := w.Write(headers) // DEBUG
	//if err != nil {
	//	panic(err)
	//}
	//
	//for _, r := range repos {
	//	class := "?"
	//	switch {
	//	case strings.HasPrefix(r, "net-"):
	//		class = "Networking"
	//	case strings.HasPrefix(r, "eventing-"):
	//		class = "Eventing"
	//	case strings.HasPrefix(r, "kn-"):
	//		class = "Client"
	//	case strings.HasPrefix(r, "sample-"):
	//		class = "Samples"
	//	}
	//	_, _ = w.Write([]byte(r + "\t" + class + "\n"))
	//}
	//return nil
}
