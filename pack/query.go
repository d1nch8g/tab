// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"os"

	"fmnx.su/core/pack/msgs"
	"fmnx.su/core/pack/pacman"
)

// Parameters that will be used to execute push command.
type QueryParameters struct {
	// List outdated packages.
	Outdated bool `short:"o" long:"outdated"`
	// Get information about package.
	Info []bool `short:"i" long:"info"`
	// List package files.
	List []bool `short:"l" long:"list"`
}

var QueryHelp = `Query packages

options:
	-i, --info     View package information (-ii for backup files)
	-l, --list     List the files owned by the queried package
	-o, --outdated List outdated packages

usage:  pack {-Q --query} [options] <(registry)/(owner)/package(s)>`

func Query(args []string, prms ...QueryParameters) error {
	p := getOptions(prms)

	if p.Outdated {
		err := pacman.SyncList(nil, pacman.SyncParameters{
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
			Stdin:   os.Stdin,
			Sudo:    true,
			Refresh: []bool{true, true},
		})
		if err != nil {
			return err
		}

		msgs.Amsg(os.Stdout, "Outdated packages")
		return pacman.Query(nil, pacman.QueryParameters{
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
			Stdin:   os.Stdin,
			Upgrade: true,
		})
	}

	return pacman.Query(args, pacman.QueryParameters{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
		Info:   p.Info,
		List:   p.List,
	})
}
