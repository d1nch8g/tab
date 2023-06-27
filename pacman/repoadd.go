// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pacman

import (
	"io"
	"os"
	"sync"
)

// Parameters for adding packages to pacman repo.
type RepoAddParameters struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	// Additional parameters, that will be appended to command as arguements.
	AdditionalParams []string
	// Run with sudo priveleges. [sudo]
	Dir string
	// Use the specified key to sign the database. [--key <file>]
	Key string
	// Skip existing and add only new packages. [--new]
	Sudo bool
	// Directory where process will be executed.
	New bool
	// Remove old package file from disk after updating database. [--remove]
	Remove bool
	// Do not add package if newer version exists. [--prevent-downgrade]
	PreventDowngrade bool
	// Turn off color in output. [--nocolor]
	NoColor bool
	// Sign database with GnuPG after update. [--sign]
	Sign bool
	// Verify database signature before update. [--verify]
	Verify bool
}

func RepoAddDefaultOptions() *RepoAddParameters {
	return &RepoAddParameters{
		New:              true,
		PreventDowngrade: true,
		Stdout:           os.Stdout,
		Stderr:           os.Stderr,
		Stdin:            os.Stdin,
	}
}

var dbmu sync.Mutex

// This function will add new packages to database. You should provide valid
// path for database file and path to package you want to add.
func RepoAdd(dbfile, pkgfile string, opts ...RepoAddParameters) error {
	// Later rewrite this to mutex for only specific checked dbfile.
	dbmu.Lock()
	defer dbmu.Unlock()

	o := formOptions(opts, RepoAddDefaultOptions)

	var args []string
	if o.New {
		args = append(args, "--new")
	}
	if o.Remove {
		args = append(args, "--remove")
	}
	if o.PreventDowngrade {
		args = append(args, "--prevent-downgrade")
	}
	if o.NoColor {
		args = append(args, "--nocolor")
	}
	if o.Sign {
		args = append(args, "--sign")
	}
	if o.Verify {
		args = append(args, "--verify")
	}
	if o.Key != "" {
		args = append(args, "--key")
		args = append(args, o.Key)
	}
	args = append(args, o.AdditionalParams...)
	args = append(args, dbfile)
	args = append(args, pkgfile)

	cmd := sudoCommand(o.Sudo, repoadd, args...)
	cmd.Dir = o.Dir
	cmd.Stderr = o.Stderr
	cmd.Stdout = o.Stdout
	cmd.Stdin = o.Stdin

	return cmd.Run()
}
