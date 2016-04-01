package main

//
// A rough sketch of a namerd controller utility.
//

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/buoyantio/namectl/namerd"
	// TODO "github.com/spf13/cobra"
)

var usageOverview = `namectl controls the namerd rpc route management service.

Usage:
  namectl [flags]
  namectl [command]

Available Commands:
  list		List all Dtab names.
  get		Get a Dtab by name.
  new		Create a new Dtab.
  update	Update an existing Dtab.
`

func exUsage(msg string) {
	fmt.Fprint(os.Stderr, usageOverview)
	fmt.Fprintf(os.Stderr, "\nerror: %s\n", msg)
	os.Exit(64)
}

func main() {
	baseURLStr := flag.String("base-url", os.Getenv("NAMECTL_BASE_URL"),
		"base url to reach namerd's ctl api")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		exUsage("command required")
	}
	cmd := args[0]
	args = args[1:]

	if *baseURLStr == "" {
		exUsage("--base-url must be specified")
	}

	baseURL, err := url.Parse(*baseURLStr)
	if err != nil {
		exUsage(fmt.Sprintf("invalid --base-url: %s: %s", *baseURLStr, err))
	}
	if baseURL.Scheme == "" || baseURL.Host == "" {
		exUsage("invalid --base-url: " + *baseURLStr)
	}
	ctl := namerd.NewHttpController(baseURL, &http.Client{})

	getDtabName := func() string {
		name := args[0]
		if name == "" {
			exUsage("empty dtab name")
		}
		return name
	}
	switch cmd {
	case "list":
		if len(args) != 0 {
			exUsage("list does not accept arguments")
		}
		names, err := ctl.List()
		dieIf(err)
		fmt.Println(strings.Join(names, "\n"))

	case "get":
		if len(args) != 1 {
			exUsage("get <name>")
		}
		name := getDtabName()
		rsp, err := ctl.Get(name)
		dieIf(err)
		fmt.Println(strings.Join(strings.SplitAfter(string(rsp.Dtab), ";"), "\n"))

	case "new":
		if len(args) != 2 {
			exUsage(fmt.Sprintf("%s <name> <file>", cmd))
		}
		name := getDtabName()
		dtab := readDtabPath(args[1])
		_, err := ctl.Create(name, dtab)
		dieIf(err)

	case "update":
		if len(args) != 2 {
			exUsage(fmt.Sprintf("%s <name> <file>", cmd))
		}
		name := getDtabName()
		dtab := readDtabPath(args[1])
		v, err := ctl.Get(name)
		dieIf(err)
		_, err = ctl.Update(name, dtab, v.Version)
		dieIf(err)

	default:
		exUsage("unexpected command: " + cmd)
	}
}

func dieIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func readDtabPath(path string) namerd.Dtab {
	var file io.Reader
	var err error
	switch path {
	case "":
		exUsage("empty dtab path")
	case "-":
		file = os.Stdin
	default:
		file, err = os.Open(path)
		dieIf(err)
	}
	dtab, err := ioutil.ReadAll(file)
	dieIf(err)
	return namerd.Dtab(dtab)
}
