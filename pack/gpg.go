// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"errors"
	"io"
	"os/exec"
	"strings"
)

// Parameters that will be used to execute gpg util commands.
type GpgParameters struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	// Export public GPG key armor
	Export bool
	// Set gpg key id as git signing key (provide as arguement)
	Gitkey bool
	// List secret keys with their IDs
	Privid bool
	// List public keys with their IDs
	Pubring bool
}

func gpgdefault() *GpgParameters {
	return &GpgParameters{}
}

// Push your package to registry.
func Gpg(args []string, prms ...GpgParameters) error {
	p := formOptions(prms, gpgdefault)

	switch {
	case p.Export:
		return Export(p)
	case p.Gitkey:
		return SetGit(args, p)
	case p.Privid:
		return Privid(p)
	case p.Pubring:
		return Pubring(p)
	}

	return errors.ErrUnsupported
}

// Export public GPG key, which can be added to gitea/gitlab/github.
func Export(p *GpgParameters) error {
	ident, err := GnuPGidentity()
	if err != nil {
		return err
	}
	gpgemail := strings.Replace(strings.Split(ident, " <")[1], ">", "", 1)
	cmd := exec.Command("gpg", "--armor", "--export", gpgemail)
	cmd.Stdout = p.Stdout
	cmd.Stderr = p.Stderr
	cmd.Stdin = p.Stdin
	return cmd.Run()
}

// Set new signing id in git.
func SetGit(args []string, p *GpgParameters) error {
	if len(args) != 1 {
		return errors.New("provide key string as arguement")
	}
	cmd := exec.Command("git", "--config", "--global", "user.signingkey", args[0])
	cmd.Stdout = p.Stdout
	cmd.Stderr = p.Stderr
	cmd.Stdin = p.Stdin
	return cmd.Run()
}

// List private gpg key IDs.
func Privid(p *GpgParameters) error {
	cmd := exec.Command("gpg", "-K")
	cmd.Stdout = p.Stdout
	cmd.Stderr = p.Stderr
	cmd.Stdin = p.Stdin
	return cmd.Run()
}

// List public gpg IDs.
func Pubring(p *GpgParameters) error {
	cmd := exec.Command("gpg", "-k")
	cmd.Stdout = p.Stdout
	cmd.Stderr = p.Stderr
	cmd.Stdin = p.Stdin
	return cmd.Run()
}
