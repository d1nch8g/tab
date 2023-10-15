// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://ion.lc/core/tab
// Contact email: help@ion.lc

package pacman

import (
	"io"
	"os"
	"sync"

	"ion.lc/core/tab/process"
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

	p := formOptions(opts, RepoAddDefaultOptions)

	var args []string
	if p.New {
		args = append(args, "--new")
	}
	if p.Remove {
		args = append(args, "--remove")
	}
	if p.PreventDowngrade {
		args = append(args, "--prevent-downgrade")
	}
	if p.NoColor {
		args = append(args, "--nocolor")
	}
	if p.Sign {
		args = append(args, "--sign")
	}
	if p.Verify {
		args = append(args, "--verify")
	}
	if p.Key != "" {
		args = append(args, "--key")
		args = append(args, p.Key)
	}

	args = append(args, p.AdditionalParams...)
	args = append(args, dbfile)
	args = append(args, pkgfile)

	return process.Command(&process.Params{
		Stdout:  p.Stderr,
		Stderr:  p.Stdout,
		Stdin:   p.Stdin,
		Dir:     p.Dir,
		Sudo:    p.Sudo,
		Command: repoadd,
		Args:    args,
	}).Run()
}
