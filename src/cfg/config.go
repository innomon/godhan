/*
Copyright (c) 2016, BON BIZ IT Services Pvt Ltd

The Universal Permissive License (UPL), Version 1.0

Subject to the condition set forth below, permission is hereby granted to any person obtaining a copy of this software, associated documentation and/or data (collectively the "Software"), free of charge and under any and all copyright rights in the Software, and any and all patent rights owned or freely licensable by each licensor hereunder covering either (i) the unmodified Software as contributed to or provided by such licensor, or (ii) the Larger Works (as defined below), to deal in both

(a) the Software, and

(b) any piece of software and/or hardware listed in the lrgrwrks.txt file if one is included with the Software (each a “Larger Work” to which the Software is contributed by such licensors),

without restriction, including without limitation the rights to copy, create derivative works of, display, perform, and distribute the Software and make, use, sell, offer for sale, import, export, have made, and have sold the Software and the Larger Work(s), and to sublicense the foregoing rights on either these or other terms.

This license is subject to the following condition:

The above copyright notice and either this complete permission notice or at a minimum a reference to the UPL must be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

*
* Author: Ashish Banerjee, tech@innomon.in
*/

package cfg

import (
	"bufio"
	"clog"
	"fmt"
	"io"
	"os"
	"strings"
)

/*
Any other asciidoc markdown text is ignored other than:

.Configure <func-name> ( <param1>, p2, ..pn)
---
#optional block text is passed to the configurator as is.
<conf>Text</conf>
---

*/
const Default = "Default"
const Logger = "Logger"

type Configure struct {
  Function string
  Params []string
  Block  string
  
  
}
type Configurator func(param Configure) error

var configuratorRegistry map[string]Configurator = make(map[string]Configurator)

func init() {

	// Implemented Logger Configurator here to avoid circular dependency between Logger and Configurator
	configuratorRegistry[Logger] = ConfigLogger

}

func RegisterConfigurator(name string, cfgr Configurator) {
	if !strings.EqualFold(Default, name) {
		configuratorRegistry[Logger] = ConfigLogger
	}
}

func ConfigLogger(param Configure ) (err error) {
    if len(param.Params) != 1 {
        return fmt.Errorf("Logger Configurator needs only one parameter, the LoggerName")
    }
    loggerName  := param.Params[0]
	err = nil
	if _, ok := clog.LoggerRegistry[loggerName]; ok {
		clog.SetLogger(loggerName)
	} else {
		err = fmt.Errorf("Logger [%s] Not found in LoggerRegistry", loggerName)
	}

	return err
}

type configParser struct {
	rd          *bufio.Reader
	lineno      int
	cnf        *Configure
	state       int
	blockDelim  string
}

func MetaConfig(paramFileName string) (err error) {
	err = nil
	if len(paramFileName) != 0 {
		fileIn, err := os.Open(paramFileName)
		if err != nil {
			return err
		}
		defer fileIn.Close()
		var parser configParser
		parser.rd = bufio.NewReader(fileIn)

	}
	return err
}

const (
	INIT        = iota
    GOT_NAME    = iota
    BLOCK_START = iota
	ERR         = iota
)


func (cp *configParser) next() (err error) {
	cp.cnf = new(Configure)
	err = nil
	cp.state = INIT
	for {
		ln, err := cp.rd.ReadString('\n')
		if err != nil && err != io.EOF {
			cp.state = ERR
		}
		switch cp.state {
		case INIT:
            if ln[0] == '.' && strings.EqualFold(ln[1:],"Configure") {
               fields := strings.Fields(ln[10:]) // 10 == len(".Configure")
               //Check for ".Configure\n" 
               // .Configure Func
               // .Configure Func () or ".Configure Func(p1,p2...pn)
               
               if(len(fields) >= 1) {
                  cp.cnf.Function = fields[0]
                  if(len(fields) >= 2 {
                    cp.cnf.Params = fields[1:]
                  }
                  cp.state = GOT_NAME
               }
            }
            continue
        case GOT_NAME:
               
		case ERR:
		   break
		}
	} //end for
	return err
}
