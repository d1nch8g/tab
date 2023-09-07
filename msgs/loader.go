// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package msgs

import (
	"fmt"
	"io"
	"strings"

	"github.com/mitchellh/ioprogress"
	"golang.org/x/term"
)

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
		prefix     = fmt.Sprintf("(%d/%d) %s", p.Current, p.Total, p.Msg)
		prefixlen  = len(prefix)
		loader     = int(float64(width) * 0.35)
		percentage = 4
	)

	switch {
	// Very slim terminal. Trimmed prefix and loading percentage are visible.
	case width < prefixlen+percentage+1:
		cutprefix := prefix[:width-percentage-4] + "..."

		return ioprogress.DrawTerminalf(p.Output, func(current, total int64) string {
			progress := float32(current) / float32(total) * 100

			return fmt.Sprintf("%s %.0f", cutprefix, progress) + "%"
		})

	// Slim terminal. Full prefix and loading percentage are visible.
	case width < prefixlen+loader+percentage+3:
		padding := strings.Repeat(" ", width-prefixlen-percentage)

		return ioprogress.DrawTerminalf(p.Output, func(current, total int64) string {
			progress := float32(current) / float32(total) * 100

			return fmt.Sprintf("%s%s%.0f", prefix, padding, progress) + "%"
		})

	// Normal size terminal. Full prefix, full loader and percetage are visible.
	default:
		padding := strings.Repeat(" ", width-prefixlen-percentage-loader-3)

		return ioprogress.DrawTerminalf(p.Output, func(current, total int64) string {
			progress := float32(current) / float32(total) * 100
			curr := int((float64(current) / float64(total)) * float64(loader))
			loader := fmt.Sprintf(
				"[%s%s]", strings.Repeat("#", curr),
				strings.Repeat("-", loader-curr),
			)

			return fmt.Sprintf("%s%s%s %.0f", prefix, padding, loader, progress) + "%"
		})
	}
}
