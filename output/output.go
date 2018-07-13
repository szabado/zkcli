/*
   Copyright 2014 Outbrain Inc.

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

// output provides with controlled printing functions
package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var Out *os.File

func init() {
	Out = os.Stdout
}

type Printer interface {
	Printf(format string, a ...interface{})
	PrintArray(array []string)
}

type TxtPrinter struct {
	OmitTrailingNL bool
}

func (p TxtPrinter) Print(data string) {
	if p.OmitTrailingNL {
		fmt.Fprint(Out, data)
	} else {
		fmt.Fprintln(Out, data)
	}
}

func (p TxtPrinter) Printf(format string, a ... interface{}) {
	fmt.Fprintf(Out, format, a...)
}

func (p TxtPrinter) PrintArray(stringArray []string) {
	s := strings.Join(stringArray, "\n")
	if p.OmitTrailingNL {
		fmt.Fprint(Out, s)
	} else {
		fmt.Fprintln(Out, s)
	}
}

type JSONPrinter struct{}

func (p JSONPrinter) Printf(format string, a ... interface{}) {
	s := fmt.Sprintf(format, a...)
	b, _ := json.Marshal(s)
	fmt.Fprintln(Out, string(b))
}

func (_ JSONPrinter) PrintArray(stringArray []string) {
	s, _ := json.Marshal(stringArray)
	fmt.Fprintln(Out, string(s))
}
