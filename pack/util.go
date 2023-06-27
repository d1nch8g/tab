// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"fmnx.su/core/pack/msgs"
)

// Parameters for util.
type UtilParameters struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	// Generate GnuPG key.
	Gen bool
	// Export existing GnuPG key armored string.
	Armor bool
	// Run gpg --recv-key to avoid pacman signing problems.
	Recv bool
	// Set packager in pacman.conf
	Setpkgr bool
	// Generate flutter template.
	Flutter bool
	// Generate go cli utility template.
	Gocli bool
}

func utildefault() *UtilParameters {
	return &UtilParameters{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
}

func Util(args []string, prms ...UtilParameters) error {
	p := formOptions(prms, utildefault)

	switch {
	case p.Recv:
		return recv(args[0], p)
	case p.Setpkgr:
		return setpkgr(args[0], p)
	case p.Gen:
		return generate(p)
	case p.Armor:
		return armored(p.Stdout)
	case p.Flutter:
		return fluttertemplate()
	case p.Gocli:
		return goclitemplate()
	}
	return errors.New("specify command options, run 'pack -Uh'")
}

// Add packager line to makepkg.conf
func setpkgr(pkgr string, p *UtilParameters) error {
	c := fmt.Sprintf("echo PACKAGER='%s' >> /etc/makepkg.conf", pkgr)
	cmd := exec.Command("bash", "-c", c)
	cmd.Stdout = p.Stdout
	cmd.Stderr = p.Stderr
	cmd.Stdin = p.Stdin
	return cmd.Run()
}

// Recieve gpg key.
func recv(id string, p *UtilParameters) error {
	cmd := exec.Command("gpg", "--recv-key", id)
	cmd.Stdout = p.Stdout
	cmd.Stderr = p.Stderr
	cmd.Stdin = p.Stdin
	return cmd.Run()
}

// Generate new GPG key with user input and etc.
func generate(p *UtilParameters) error {
	cmd := exec.Command("gpg", "--armor", "--export")
	cmd.Stdout = p.Stdout
	cmd.Stderr = p.Stderr
	cmd.Stdin = p.Stdin
	return cmd.Run()
}

// Return armored public key string from GnuPG.
func armored(o io.Writer) error {
	cmd := exec.Command("gpg", "--armor", "--export")
	cmd.Stdout = o
	return call(cmd)
}

// Function generates project template for flutter desktop application based on
// current directory name and identity in GnuPG.
func fluttertemplate() error {
	ident, err := gnuPGIdentity()
	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	splt := strings.Split(dir, "/")
	n := splt[len(splt)-1]

	d := fmt.Sprintf(msgs.Desktop, n, n, n, n, n)
	derr := os.WriteFile(n+".desktop", []byte(d), 0600)

	s := fmt.Sprintf(msgs.ShFile, n, n)
	serr := os.WriteFile(n+".sh", []byte(s), 0600)

	p := fmt.Sprintf(msgs.PKGBUILDflutter, ident, n, n, n, n, n, n, n, n)
	perr := os.WriteFile(`PKGBUILD`, []byte(p), 0600)

	return errors.Join(derr, serr, perr)
}

// Function generates project template for go cli utility based on
// current directory name and identity in GnuPG.
func goclitemplate() error {
	ident, err := gnuPGIdentity()
	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	splt := strings.Split(dir, "/")
	n := splt[len(splt)-1]

	p := fmt.Sprintf(msgs.PKGBUILDgocli, ident, n, n)
	return os.WriteFile(`PKGBUILD`, []byte(p), 0600)
}
