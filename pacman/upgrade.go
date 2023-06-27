// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pacman

import (
	"io"
	"os"
	"strings"
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
	o := formOptions(opts, UpgradeDefault)

	args := []string{"-U"}
	if o.Needed {
		args = append(args, "--needed")
	}
	if o.NoConfirm {
		args = append(args, "--noconfirm")
	}
	if o.NoProgressBar {
		args = append(args, "--noprogressbar")
	}
	if o.NoScriptlet {
		args = append(args, "--noscriptlet")
	}
	if o.AsDeps {
		args = append(args, "--asdeps")
	}
	if o.AsExplict {
		args = append(args, "--asexplict")
	}
	args = append(args, o.AdditionalParams...)
	args = append(args, files...)

	cmd := sudoCommand(o.Sudo, pacman, args...)
	cmd.Stdout = o.Stdout
	cmd.Stderr = o.Stderr
	cmd.Stdin = o.Stdin

	mu.Lock()
	defer mu.Unlock()
	return cmd.Run()
}
