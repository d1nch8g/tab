// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"fmnx.su/core/pack/msgs"
)

// Parameters to generate PKGBUILD templates.
type TmplParameters struct {
	// Generate default PKGBUILD tempalte
	Default bool `short:"t" long:"default"`
	// Templtate for flutter project
	Flutter bool `long:"flutter"`
	// Templtate for CLI tool in go
	Gocli bool `long:"gocli"`
}

var TmplHelp = `Template operations

options:
	-t, --default Default PKGBUILD template /usr/share/pacman/PKGBUILD.proto
	    --flutter Templtate for flutter project 
	    --gocli   Templtate for CLI tool in go

usage:  pack {-G --gpg} [options] <(args)>`

func Tmpl(args []string, prms ...TmplParameters) error {
	p := getOptions(prms)

	switch {
	case p.Default:
		return GenWrap(p, DefaultTemplate)
	case p.Flutter:
		return GenWrap(p, FlutterTemplate)
	case p.Flutter:
		return GenWrap(p, GoCliTemplate)
	}

	return errors.ErrUnsupported
}

// Wrapper function to inform that template is generated succesfully.
func GenWrap(p *TmplParameters, f func() error) error {
	err := f()
	if err != nil {
		return err
	}
	msgs.Amsg(os.Stdout, "Template generated")
	return nil
}

// Generate default project template.
func DefaultTemplate() error {
	b, err := os.ReadFile("/usr/share/pacman/PKGBUILD.proto")
	if err != nil {
		return err
	}
	return os.WriteFile("PKGBUILD", b, os.ModePerm)
}

// Function generates project template for flutter desktop application based on
// current directory name and identity in GnuPG.
func FlutterTemplate() error {
	ident, err := GnuPGidentity()
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
func GoCliTemplate() error {
	ident, err := GnuPGidentity()
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
