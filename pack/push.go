// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"fmnx.su/core/pack/creds"
	"fmnx.su/core/pack/msgs"
	"github.com/mitchellh/ioprogress"
)

// Parameters that will be used to execute push command.
type PushParameters struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	// Directory to read package files and signatures.
	Directory string `short:"d" long:"dir" default:"/var/cache/pacman/pkg"`
	// Which protocol to use for connection.
	Insecure bool `short:"i" long:"insecure"`
	// Custom distribution for which package is built.
	Distro string `short:"s" long:"distro" default:"archlinux"`
}

var PushHelp = `Push cached packages

options:
	-d, --dir <dir> Use custom source dir with packages (default pacman cache)
	-i, --insecure  Push package over HTTP instead of HTTPS
	-s, --distro    Assign custom distribution in registry (default archlinux)
	-i, --insecure  Use HTTP protocol for API call

usage:  pack {-P --push} [options] <registry/owner/package(s)>`

// Push your package to registry.
func Push(args []string, prms ...PushParameters) error {
	p := getOptions(prms)

	msgs.Amsg(p.Stdout, "Preparing pushed packages")

	cachedpkgs, err := listPkgFilenames(p.Directory)
	if err != nil {
		return err
	}
	msgs.Smsg(p.Stdout, "Scanning cached packages", 2, 3)

	mds, err := prepareMetadata(p.Directory, cachedpkgs, args)
	if err != nil {
		return err
	}
	msgs.Smsg(p.Stdout, "Preparing package metadata", 3, 3)

	msgs.Amsg(p.Stdout, "Pushing packages")
	for i, md := range mds {
		err = push(*p, md, i+1, len(mds))
		if err != nil {
			return err
		}
	}
	return nil
}

type PackageMetadata struct {
	Name     string
	FileName string
	Addr     string
	Owner    string
}

// Collect metadata about packages, ensure all packages could be pushed.
func prepareMetadata(dir string, filenames, pkgs []string) ([]PackageMetadata, error) {
	var mds []PackageMetadata
	for _, pkg := range pkgs {
		var (
			name    string
			owner   string
			address string
		)

		splt := strings.Split(pkg, "/")
		switch len(splt) {
		case 1:
			return nil, errors.New("no registry to push: " + pkg)
		case 2:
			address = splt[0]
			name = splt[1]
		case 3:
			address = splt[0]
			owner = splt[1]
			name = splt[2]
		}

		filenames, err := FilterFilenames(filenames, name)
		if err != nil {
			return nil, err
		}
		for _, filename := range filenames {
			mds = append(mds, PackageMetadata{
				Name:     name,
				FileName: filename,
				Addr:     address,
				Owner:    owner,
			})
		}
	}
	return mds, nil
}

// Filter filenames related to required package.
func FilterFilenames(filenames []string, pkg string) ([]string, error) {
	var rez []string
	for _, filename := range filenames {
		pkgsplt := strings.Split(filename, "-")
		if len(pkgsplt) < 4 {
			return nil, errors.New("not valid package file name: " + filename)
		}
		if !(strings.Join(pkgsplt[:len(pkgsplt)-3], "-") == pkg) {
			continue
		}
		rez = append(rez, filename)
	}
	return rez, nil
}

// List file names in provided cache directory.
func listPkgFilenames(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.New("unable to get directory contents: " + dir)
	}
	var fns []string
	for _, direntry := range entries {
		filename := direntry.Name()
		if strings.HasSuffix(filename, ".pkg.tar.zst") {
			fns = append(fns, filename)
		}
	}
	return fns, nil
}

// This function pushes package to registry via http/https.
func push(pp PushParameters, m PackageMetadata, i, t int) error {
	pkgpath := path.Join(pp.Directory, m.FileName)
	packagefile, err := os.Open(pkgpath)
	if err != nil {
		return err
	}
	pkgInfo, err := os.Stat(pkgpath)
	if err != nil {
		return err
	}

	pkgsign, err := os.ReadFile(pkgpath + ".sig")
	if err != nil {
		return err
	}

	protocol := "https"
	if pp.Insecure {
		protocol = "http"
	}

	req, err := http.NewRequest(
		http.MethodPut,
		protocol+"://"+path.Join(
			m.Addr, "api/packages", m.Owner, "arch/push",
			m.FileName, pp.Distro, hex.EncodeToString(pkgsign),
		),
		&ioprogress.Reader{
			Reader: packagefile,
			Size:   pkgInfo.Size(),
			DrawFunc: msgs.Loader(&msgs.LoaderParameters{
				Current: i,
				Total:   t,
				Msg: fmt.Sprintf(
					"Pushing %s to %s...", m.FileName,
					path.Join(m.Addr, m.Owner),
				),
				Output: pp.Stdout,
			}),
		},
	)
	if err != nil {
		return err
	}

	login, pass, err := creds.Get(protocol, m.Addr)
	if err != nil {
		login, pass, err = creds.Create(protocol, m.Addr, pp.Stdin, pp.Stdout)
		if err != nil {
			return err
		}
	}

	req.SetBasicAuth(login, pass)

	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Join(err, errors.New(resp.Status))
		}
		return fmt.Errorf("%s %s", resp.Status, string(b))
	}
	return nil
}
