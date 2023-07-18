// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package sudo

import "os/exec"

const (
	sudo = `sudo`
	doas = `doas`
)

func Command(elevate bool, command string, args ...string) *exec.Cmd {
	_, err := exec.LookPath(sudo)
	if elevate && err == nil {
		args = append([]string{command}, args...)
		return exec.Command(sudo, args...)
	}

	_, err = exec.LookPath(doas)
	if elevate && err == nil {
		args = append([]string{command}, args...)
		return exec.Command(doas, args...)
	}

	return exec.Command(command, args...)
}
