// modeled upon cmd/cayley.go

package main

import (
	"flag"
	"fmt"
	"gopkg.in/xmlpath.v2"
	"log"
	"os"
)

var (
	xml   = flag.String("xml", "input.xml", "xml formatted input file name")
	xpath = flag.String("xp", "", "xpath expression")
)

// Filled in by `go build ldflags="-X main.Version `ver`"`.
var (
	BuildDate string
	Version   string
)

func usage() {
	fmt.Fprintln(os.Stderr, "syntax:\n <commmand> ...options...\n")
	for key, cmd := range Cmds {
		fmt.Fprintf(os.Stderr, "%s: %s\n", key, cmd.Usage)
	}
	fmt.Fprintf(os.Stderr, "\n options:\n")
	flag.PrintDefaults()
}

type Command func()

type Commander struct {
	Exec  Command
	Usage string
}

var Cmds map[string]Commander = make(map[string]Commander)

func init() {
	flag.Usage = usage
	Cmds["help"] = Commander{Exec: usage, Usage: "Display help"}
	Cmds["dump"] = Commander{Exec: xmldump, Usage: "dumps the xml specified by option -xml"}
	Cmds["xpath"] = Commander{Exec: dumpPath, Usage: "prints the node specified by option -xp from file -xml <file name>"}

}

func main() {

	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "xpath test")
		usage()
		os.Exit(1)
	}
	cmd := os.Args[1]
	os.Args = append(os.Args[:1], os.Args[2:]...)
	flag.Parse()

	if cmdr, ok := Cmds[cmd]; ok {
		cmdr.Exec()
	} else {
		usage()
		os.Exit(1)
	}

}

// dumps the xpath Node after parsing the file specified by option -xml
func dumpPath() {
	xpth := *xpath

	if xpth == "" {
		fmt.Fprintln(os.Stderr, "XPath expression not specified, please use -xp option")
		return
	}
	path := xmlpath.MustCompile(xpth)

	file, err := os.Open(*xml) // For read access.
	if err != nil {
		log.Fatal(err)
	}
	root, err := xmlpath.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	if value, ok := path.String(root); ok {
		fmt.Fprintln(os.Stderr, "Found:", value)
	} else {
		fmt.Fprintln(os.Stderr, "Not found: ", xpth)
	}

}

// dumps the root Node after parsing the file specified by option -xml
func xmldump() {
	fmt.Fprintln(os.Stderr, *xml)
	file, err := os.Open(*xml) // For read access.
	if err != nil {
		log.Fatal(err)
	}
	root, err := xmlpath.Parse(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(os.Stderr, root.String())
}
