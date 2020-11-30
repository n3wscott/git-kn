package knative

import (
	"io"
	"strings"

	"github.com/n3wscott/git-kn/pkg/printers"
)

type Repos struct {
	WG    string
	Class string
	Name  string
}

var headers = []byte("NAME\tKIND\n")

func ClassifyRepos(repos []string, output io.Writer) []Repos {

	w := printers.GetNewTabWriter(output) // DEBUG
	defer w.Flush()
	_, err := w.Write(headers) // DEBUG
	if err != nil {
		panic(err)
	}

	for _, r := range repos {
		class := "?"
		switch {
		case strings.HasPrefix(r, "net-"):
			class = "Networking"
		case strings.HasPrefix(r, "eventing-"):
			class = "Eventing"
		case strings.HasPrefix(r, "kn-"):
			class = "Client"
		case strings.HasPrefix(r, "sample-"):
			class = "Samples"
		}
		_, _ = w.Write([]byte(r + "\t" + class + "\n"))
	}
	return nil
}
