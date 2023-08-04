// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"io"
	"os"

	"fmnx.su/core/pack/msgs"
	"fmnx.su/core/pack/pacman"
)

// Parameters that will be used to execute push command.
type QueryParameters struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	// List outdated packages.
	Outdated bool
	// Get information about package.
	Info []bool
	// List package files.
	List []bool
}

func querydefault() *QueryParameters {
	return &QueryParameters{}
}

func Query(args []string, prms ...QueryParameters) error {
	p := formOptions(prms, querydefault)

	if p.Outdated {
		err := pacman.SyncList(nil, pacman.SyncParameters{
			Stdout:  p.Stdout,
			Stderr:  p.Stderr,
			Stdin:   p.Stdin,
			Sudo:    true,
			Refresh: []bool{true, true},
		})
		if err != nil {
			return err
		}

		msgs.Amsg(p.Stdout, "Outdated packages")
		return pacman.Query(nil, pacman.QueryParameters{
			Stdout:  p.Stdout,
			Stderr:  p.Stderr,
			Stdin:   p.Stdin,
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
