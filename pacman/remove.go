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

// Optional parameters for pacman remove command.
type RemoveParameters struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	// Run with sudo priveleges. [sudo]
	Sudo bool
	// Do not ask for any confirmation. [--noconfirm]
	NoConfirm bool
	// Remove with all unnecessary packages. [--recursive]
	Recursive bool
	// Remove with all explicitly installed deps. [-ss]
	ForceRecursive bool
	// Use cascade when removing packages. [--cascade]
	Cascade bool
	// Remove configuration files aswell. [--nosave]
	WithConfigs bool
	// Additional parameters, that will be appended to command as arguements.
	AdditionalParams []string
}

func RemoveDefault() *RemoveParameters {
	return &RemoveParameters{
		Recursive:   true,
		WithConfigs: true,
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		Stdin:       os.Stdin,
	}
}

// Remove packages from system.
func Remove(pkgs string, opts ...RemoveParameters) error {
	return RemoveList(strings.Split(pkgs, " "), opts...)
}

// Remove packages from system.
func RemoveList(pkgs []string, opts ...RemoveParameters) error {
	o := formOptions(opts, RemoveDefault)

	var args = []string{"-R"}
	if o.NoConfirm {
		args = append(args, "--noconfirm")
	}
	if o.Recursive {
		args = append(args, "--recursive")
	}
	if o.ForceRecursive {
		args = append(args, "-ss")
	}
	if o.Cascade {
		args = append(args, "--cascade")
	}
	if o.WithConfigs {
		args = append(args, "--nosave")
	}
	args = append(args, o.AdditionalParams...)
	args = append(args, pkgs...)

	cmd := sudoCommand(o.Sudo, pacman, args...)
	cmd.Stdout = o.Stdout
	cmd.Stderr = o.Stderr
	cmd.Stdin = o.Stdin

	mu.Lock()
	defer mu.Unlock()
	return cmd.Run()
}
