// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://ion.lc/core/tab
// Contact email: help@ion.lc

package msgs

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
)

// Write an announcement message with dots prefix and bold text to provided
// io.Writer.
func Amsg(w io.Writer, msg string) {
	dots := color.New(color.FgWhite, color.Bold, color.FgHiBlue).Sprintf(":: ")
	msg = color.New(color.Bold).Sprintf(msg)
	w.Write([]byte(dots + msg + "...\n"))
}

// Write step message, with enumeration which should represent state of program
// execution.
func Smsg(w io.Writer, msg string, i, t int) {
	w.Write([]byte(fmt.Sprintf("(%d/%d) %s...\n", i, t, msg)))
}

// Request an input from user.
func Inp(msg string, w io.Writer, r io.Reader, hidden bool) (string, error) {
	w.Write([]byte(msg))
	if hidden {
		w.Write([]byte("\033[8m"))
		defer w.Write([]byte("\033[28m"))
	}
	reader := bufio.NewReader(r)
	s, err := reader.ReadString('\n')
	if err != nil {
		return ``, err
	}
	return strings.ReplaceAll(s, "\n", ""), nil
}
