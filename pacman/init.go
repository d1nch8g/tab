// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://ion.lc/core/tab
// Contact email: help@ion.lc

package pacman

import (
	"sync"
)

// Used commands.
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
