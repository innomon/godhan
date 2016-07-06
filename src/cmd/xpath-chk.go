/*
* Copyright (c) 2016, BON BIZ IT Services Pvt LTD.
*
* The Universal Permissive License (UPL), Version 1.0
*
* Subject to the condition set forth below, permission is hereby granted to any person obtaining a copy of this software, associated documentation and/or data (collectively the "Software"), free of charge and under any and all copyright rights in the Software, and any and all patent rights owned or freely licensable by each licensor hereunder covering either (i) the unmodified Software as contributed to or provided by such licensor, or (ii) the Larger Works (as defined below), to deal in both

* (a) the Software, and

* (b) any piece of software and/or hardware listed in the lrgrwrks.txt file if one is included with the Software (each a “Larger Work” to which the Software is contributed by such licensors),
*
* without restriction, including without limitation the rights to copy, create derivative works of, display, perform, and distribute the Software and make, use, sell, offer for sale, import, export, have made, and have sold the Software and the Larger Work(s), and to sublicense the foregoing rights on either these or other terms.
*
* This license is subject to the following condition:
*
* The above copyright notice and either this complete permission notice or at a minimum a reference to the UPL must be included in all copies or substantial portions of the Software.
*
* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*
* Author: Ashish Banerjee, tech@innomon.in
*/


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
