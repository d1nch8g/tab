// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pacman

import (
	"os/exec"
	"sync"
)

// Dependecy packages.
const (
	pacman  = `pacman`
	makepkg = `makepkg`
	repoadd = `repo-add`
)

// Global lock for operations with pacman database.
var mu sync.Mutex

func formOptions[Options any](arr []Options, getdefault func() *Options) *Options {
	if len(arr) != 1 {
		return getdefault()
	}
	return &arr[0]
}

func sudoCommand(sudo bool, command string, args ...string) *exec.Cmd {
	if sudo {
		args = append([]string{command}, args...)
		return exec.Command(`sudo`, args...)
	}
	return exec.Command(command, args...)
}
