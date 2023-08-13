// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

// package pack

// import (
// 	"errors"
// 	"fmt"
// 	"io"
// 	"os"
// 	"strings"

// 	"fmnx.su/core/pack/msgs"
// )

// // Parameters for util.
// type AssistParameters struct {
// 	Stdout io.Writer
// 	Stderr io.Writer
// 	Stdin  io.Reader

// 	// Export existing GnuPG key armored string.
// 	Export bool
// 	// Check compatability of identities across git, gpg and makepkg.
// 	Fix bool
// 	// Generate flutter template.
// 	Flutter bool
// 	// Generate go cli utility template.
// 	Gocli bool
// }

// func utildefault() *AssistParameters {
// 	return &AssistParameters{
// 		Stdout: os.Stdout,
// 		Stderr: os.Stderr,
// 		Stdin:  os.Stdin,
// 	}
// }

// func Assist(args []string, prms ...AssistParameters) error {
// 	p := formOptions(prms, utildefault)

// 	switch {
// 	case p.Export:
// 		return Export(p.Stdout, p.Stderr)
// 	case p.Fix:
// 		return Fix()
// 	case p.Flutter:
// 		return FlutterTemplate()
// 	case p.Gocli:
// 		return GoCliTemplate()
// 	}

// 	return errors.New("specify assist option, run 'pack -Ah'")
// }

// // Check/Fix compatability of identities in git, gpg and makepkg.
// func Fix() error {
// 	err := ValidatePackager()
// 	if err != nil {
// 		return err
// 	}

// 	err = ValidateGitUser()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func ValidateGitUser() error {
// 	gitname, err := callOut("git", "config", "user.name")
// 	if err != nil {
// 		return fmt.Errorf("unable to get git name")
// 	}

// 	gitemail, err := callOut("git", "config", "user.email")
// 	if err != nil {
// 		return fmt.Errorf("unable to get git email")
// 	}

// 	gitidentity := gitname + " <" + gitemail + ">"

// 	gpgidentity, err := GnuPGidentity()
// 	if err != nil {
// 		return err
// 	}

// 	if gpgidentity != gitidentity {
// 		return fmt.Errorf(msgs.ErrGitUserMissmatch, gpgidentity, gitidentity)
// 	}

// 	gitsignkey, err := callOut("git", "config", "user.signingkey")
// 	if err != nil {
// 		return fmt.Errorf("unable to get git signingkey")
// 	}

// 	gpginfo, err := callOut("gpg", "-K")
// 	if err != nil {
// 		return err
// 	}

// 	if !strings.Contains(gpginfo, gitsignkey) {
// 		return fmt.Errorf(msgs.ErrGitSignKeyMissmatch, gitsignkey)
// 	}

// 	return nil
// }

// // Function generates project template for flutter desktop application based on
// // current directory name and identity in GnuPG.
// func FlutterTemplate() error {
// 	ident, err := GnuPGidentity()
// 	if err != nil {
// 		return err
// 	}

// 	dir, err := os.Getwd()
// 	if err != nil {
// 		return err
// 	}
// 	splt := strings.Split(dir, "/")
// 	n := splt[len(splt)-1]

// 	d := fmt.Sprintf(msgs.Desktop, n, n, n, n, n)
// 	derr := os.WriteFile(n+".desktop", []byte(d), 0600)

// 	s := fmt.Sprintf(msgs.ShFile, n, n)
// 	serr := os.WriteFile(n+".sh", []byte(s), 0600)

// 	p := fmt.Sprintf(msgs.PKGBUILDflutter, ident, n, n, n, n, n, n, n, n)
// 	perr := os.WriteFile(`PKGBUILD`, []byte(p), 0600)

// 	return errors.Join(derr, serr, perr)
// }

// // Function generates project template for go cli utility based on
// // current directory name and identity in GnuPG.
// func GoCliTemplate() error {
// 	ident, err := GnuPGidentity()
// 	if err != nil {
// 		return err
// 	}

// 	dir, err := os.Getwd()
// 	if err != nil {
// 		return err
// 	}
// 	splt := strings.Split(dir, "/")
// 	n := splt[len(splt)-1]

// 	p := fmt.Sprintf(msgs.PKGBUILDgocli, ident, n, n)
// 	return os.WriteFile(`PKGBUILD`, []byte(p), 0600)
// }
