// Use of this code is governed by GNU General Public License.
// Official web page: https://ion.lc/core/tab
// Contact email: help@ion.lc

package process

import (
	"io"
	"os"
	"os/exec"
)

const (
	sudo = `sudo`
	doas = `doas`
)

type Params struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	Sudo    bool
	Command string
	Args    []string
	Dir     string
}

func Command(p *Params) *exec.Cmd {
	var cmd *exec.Cmd

	_, err := exec.LookPath(sudo)
	if p.Sudo && err == nil {
		p.Args = append([]string{p.Command}, p.Args...)
		cmd = exec.Command(sudo, p.Args...)
	}

	_, err = exec.LookPath(doas)
	if p.Sudo && err == nil {
		p.Args = append([]string{p.Command}, p.Args...)
		cmd = exec.Command(doas, p.Args...)
	}

	if cmd == nil {
		cmd = exec.Command(p.Command, p.Args...)
	}

	if p.Stderr != nil {
		cmd.Stderr = p.Stderr
	} else {
		cmd.Stderr = os.Stderr
	}

	if p.Stdout != nil {
		cmd.Stdout = p.Stdout
	} else {
		cmd.Stdout = os.Stdout
	}

	if p.Stdin != nil {
		cmd.Stdin = p.Stdin
	} else {
		cmd.Stdin = os.Stdin
	}

	cmd.Dir = p.Dir

	return cmd
}
