// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// Parameters that will be used to execute gpg util commands.
type GpgParameters struct {
	// Export public GPG key armor
	Export bool `short:"e" long:"export"`
	// Set gpg key id as git signing key (provide as arguement)
	Gitkey bool `short:"g" long:"gitid"`
	// List secret keys with their IDs
	Privid bool `short:"p" long:"privid"`
	// List public keys with their IDs
	Pubring bool `short:"r" long:"pubring"`
}

var GpgHelp = `GPG operations

options:
	-e, --export  Export public GPG key armor
	-g, --git     Set gpg key id as git signing key (provide as arguement)
	-p, --privid  List secret keys with their IDs
	-r, --pubring List public keys with their IDs

usage:  pack {-G --gpg} [options] <(args)>`

// Push your package to registry.
func Gpg(args []string, prms ...GpgParameters) error {
	p := getOptions(prms)

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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Set new signing id in git.
func SetGit(args []string, p *GpgParameters) error {
	if len(args) != 1 {
		return errors.New("provide key string as arguement")
	}
	cmd := exec.Command("git", "--config", "--global", "user.signingkey", args[0])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// List private gpg key IDs.
func Privid(p *GpgParameters) error {
	cmd := exec.Command("gpg", "-K")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// List public gpg IDs.
func Pubring(p *GpgParameters) error {
	cmd := exec.Command("gpg", "-k")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
