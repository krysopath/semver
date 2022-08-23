package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	ver "github.com/krysopath/semver/versions"
)

var (
	emitFormat  = flag.String("format", "json", "output format on stdout")
	releaseType = flag.String("release", "", "specify a release type to increment the version: major|minor|patch")
)

func outputSingle(data string) {
	sem := ver.SemanticVersion{data}

	if len(*releaseType) > 0 {
		switch *releaseType {
		case "major":
			sem = sem.Release("major")
		case "minor":
			sem = sem.Release("minor")
		case "patch":
			sem = sem.Release("patch")
		default:
			panic("omg")
		}
	}

	var out []byte

	switch *emitFormat {
	case "json":
		out, _ = json.Marshal(sem)
	case "eval":
		out, _ = sem.MarshalEVAL()
	}

	fmt.Fprintln(
		os.Stdout,
		strings.TrimSpace(string(out)),
	)
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

func init() {
	flag.Parse()
}

func main() {
	data := input()
	outputSingle(data[0])
}
