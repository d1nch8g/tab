// Use of this code is governed by GNU General Public License.
// Official web page: https://ion.lc/core/tab
// Contact email: help@ion.lc

package pacman

import (
	"io"
	"os"
	"strings"

	"ion.lc/core/tab/process"
)

// Options to apply when searching for some package.
type UpgradeParameters struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	// Run with sudo priveleges. [sudo]
	Sudo bool
	// Do not reinstall up to date packages. [--needed]
	Needed bool
	// Do not ask for any confirmation. [--noconfirm]
	NoConfirm bool
	// Do not show a progress bar when downloading files. [--noprogressbar]
	NoProgressBar bool
	// Do not execute the install scriptlet if one exists. [--noscriptlet]
	NoScriptlet bool
	// Install packages as non-explicitly installed. [--asdeps]
	AsDeps bool
	// Install packages as explictly installed. [--asexplict]
	AsExplict bool
	// Additional parameters, that will be appended to command as arguements.
	AdditionalParams []string
}

func UpgradeDefault() *UpgradeParameters {
	return &UpgradeParameters{
		Needed:    true,
		NoConfirm: true,
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		Stdin:     os.Stdin,
	}
}

// Install packages from files.
func Upgrade(files string, opts ...UpgradeParameters) error {
	return UpgradeList(strings.Split(files, " "), opts...)
}

// Install packages from files.
func UpgradeList(files []string, opts ...UpgradeParameters) error {
	p := formOptions(opts, UpgradeDefault)

	args := []string{"-U"}
	if p.Needed {
		args = append(args, "--needed")
	}
	if p.NoConfirm {
		args = append(args, "--noconfirm")
	}
	if p.NoProgressBar {
		args = append(args, "--noprogressbar")
	}
	if p.NoScriptlet {
		args = append(args, "--noscriptlet")
	}
	if p.AsDeps {
		args = append(args, "--asdeps")
	}
	if p.AsExplict {
		args = append(args, "--asexplict")
	}
	args = append(args, p.AdditionalParams...)
	args = append(args, files...)

	mu.Lock()
	defer mu.Unlock()

	return process.Command(&process.Params{
		Stdout:  p.Stdout,
		Stderr:  p.Stderr,
		Stdin:   p.Stdin,
		Sudo:    p.Sudo,
		Command: pacman,
		Args:    args,
	}).Run()
}
