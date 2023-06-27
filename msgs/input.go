// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package msgs

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Asks the user for confirmation. A user must type in "yes" or "no" and then
// press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all
// count as confirmations. If the input is not recognized, it will ask again.
// The function does not return until it gets a valid response from the user.
func AskForConfirmation(in io.Reader, out io.Writer, msg string) bool {
	reader := bufio.NewReader(os.Stdin)

	dots := color.New(color.FgWhite, color.Bold, color.FgHiBlue).Sprintf(":: ")
	msg = color.New(color.Bold).Sprintf(msg + "? [Y/n] ")
	msg = dots + msg

	for {
		out.Write([]byte(msg))

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
