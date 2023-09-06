// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package msgs

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/mitchellh/ioprogress"
	"golang.org/x/term"
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

type LoaderParameters struct {
	Current int
	Total   int
	Msg     string
	Output  io.Writer
}

// Function that will give terminal drawer for provided message, that can be
// further used in different IO operations.
func Loader(p *LoaderParameters) func(int64, int64) error {
	width, _, err := term.GetSize(0)
	if err != nil {
		return nil
	}

	var (
		prefix       = fmt.Sprintf("(%d/%d) %s", p.Current, p.Total, p.Msg)
		loaderwidth  = int(float64(width) * 0.35)
		paddingWidth = width - len(prefix) - loaderwidth - 7
	)

	if len(prefix) > width {
		return ioprogress.DrawTerminalf(p.Output, func(progress, total int64) string {
			prcntg := float32(progress) / float32(total) * 100

			return fmt.Sprintf("%s %.0f", prefix[:width-8]+"...", prcntg) + "%"
		})
	}

	if paddingWidth+5 < 0 {
		return ioprogress.DrawTerminalf(p.Output, func(progress, total int64) string {
			prcntg := float32(progress) / float32(total) * 100

			return fmt.Sprintf("%s %.0f", prefix[:width-8]+"...", prcntg) + "%"
			// return prefix
		})
	}

	padding := strings.Repeat(" ", paddingWidth)
	return ioprogress.DrawTerminalf(p.Output, func(progress, total int64) string {
		prcntg := float32(progress) / float32(total) * 100

		current := int((float64(progress) / float64(total)) * float64(loaderwidth))
		bar := fmt.Sprintf(
			"[%s%s]", strings.Repeat("#", int(current)),
			strings.Repeat("-", int(loaderwidth-current)),
		)

		return fmt.Sprintf("%s%s%s %.0f", prefix, padding, bar, prcntg) + "%"
	})
}
