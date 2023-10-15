// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://ion.lc/core/tab
// Contact email: help@ion.lc

package msgs

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func init() {
	b, err := os.ReadFile("/etc/pacman.conf")
	if err != nil {
		fmt.Println("unable to read pacman configuration")
		os.Exit(1)
	}
	if !strings.Contains(string(b), "\nColor\n") {
		color.NoColor = true
	}
}
