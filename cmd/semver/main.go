package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/alessio/shellescape"
	ver "github.com/krysopath/semver/versions"
)

var (
	isSorting  = flag.Bool("sort", false, "sort strings from $* inputs with semver algo")
	emitFormat = flag.String("format", "json", "output format on stdout")
)

func outputSingle(data string) {
	sem := ver.SemanticVersion{data}

	var out []byte

	switch *emitFormat {
	case "json":
		out, _ = json.Marshal(sem)
	case "eval":
		out, _ = sem.MarshalEVAL()
	}

	fmt.Fprintln(
		os.Stdout,
		strings.TrimSpace(shellescape.Quote(string(out))),
	)
}

func outputSorted(data []string) {
	var sorted []ver.SemanticVersion
	for _, v := range data {
		sorted = append(sorted, ver.SemanticVersion{v})
	}
	sort.Sort(ver.ByVersion(sorted))
	out, _ := json.Marshal(sorted)
	fmt.Fprintln(os.Stdout, string(out))
}

func outputOrdered(data []string) {
	var sorted []ver.SemanticVersion
	for _, v := range data {
		sorted = append(sorted, ver.SemanticVersion{v})
	}
	out, _ := json.Marshal(sorted)
	fmt.Fprintln(os.Stdout, string(out))
}

func input() []string {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	var data string

	if fi.Mode()&os.ModeNamedPipe == 0 {
		data = os.Getenv("CI_COMMIT_TAG")
	} else {
		reader := bufio.NewReader(os.Stdin)
		data, _ = reader.ReadString('\n')
	}
	trimmed := strings.Fields(data)

	return trimmed
}

func main() {
	flag.Parse()
	data := input()
	if len(data) > 2 {
		if *isSorting {
			outputSorted(data)
		} else {
			outputOrdered(data)
		}
	} else if len(data) == 1 {
		outputSingle(data[0])
	} else {
		flag.Usage()
	}
}
