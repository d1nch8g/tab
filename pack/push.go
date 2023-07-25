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
	Directory string
	// Which protocol to use for connection.
	Insecure bool
	// Custom distribution for which package is built.
	Distro string
}

func pushdefault() *PushParameters {
	return &PushParameters{
		Directory: "/var/cache/pacman/pkg",
		Distro:    "archlinux",
	}
}

// Push your package to registry.
func Push(args []string, prms ...PushParameters) error {
	p := formOptions(prms, pushdefault)

	msgs.Amsg(p.Stdout, "Preparing pushed packages")

	email, err := gnupgEmail()
	if err != nil {
		return err
	}
	msgs.Smsg(p.Stdout, "Pushing as: "+email, 1, 3)

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
		err = push(*p, md, email, i+1, len(mds))
		if err != nil {
			return err
		}
	}
	return nil
}

// This function will be used to get email from user's GnuPG identitry.
func gnupgEmail() (string, error) {
	gnupgident, err := GnuPGidentity()
	if err != nil {
		return ``, err
	}
	return strings.ReplaceAll(strings.Split(gnupgident, "<")[1], ">", ""), nil
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
		md := PackageMetadata{}

		splt := strings.Split(pkg, "/")
		switch len(splt) {
		case 1:
			return nil, errors.New("no registry to push: " + pkg)
		case 2:
			md.Addr = splt[0]
			md.Name = splt[1]
		case 3:
			md.Addr = splt[0]
			md.Owner = splt[1]
			md.Name = splt[2]
		}

		var err error
		md.FileName, err = getLastverCachedPkgFile(md.Name, filenames)
		if err != nil {
			return nil, err
		}

		mds = append(mds, md)
	}
	return mds, nil
}

// Get lastet package from list based on package name.
func getLastverCachedPkgFile(pkg string, files []string) (string, error) {
	for i := len(files) - 1; i >= 0; i-- {
		filename := files[i]
		if !strings.HasPrefix(filename, pkg) {
			continue
		}
		pkgsplt := strings.Split(filename, "-")
		if len(pkgsplt) < 4 {
			return ``, errors.New("not valid package file name: " + filename)
		}
		if !(strings.Join(pkgsplt[:len(pkgsplt)-3], "-") == pkg) {
			continue
		}
		return filename, nil
	}
	return ``, errors.New("cannot find package in cache: " + pkg)
}

// List file names in provided cache directory.
func listPkgFilenames(dir string) ([]string, error) {
	des, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.New("unable to get directory contents: " + dir)
	}
	var fns []string
	for _, de := range des {
		fn := de.Name()
		if strings.HasSuffix(fn, ".pkg.tar.zst") {
			fns = append(fns, fn)
		}
	}
	return fns, nil
}

// This function pushes package to registry via http/https.
func push(pp PushParameters, m PackageMetadata, email string, i, t int) error {
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
		protocol+"://"+path.Join(m.Addr, "api/packages", m.Owner, "arch/push"),
		&ioprogress.Reader{
			Reader: packagefile,
			Size:   pkgInfo.Size(),
			DrawFunc: msgs.Loader(&msgs.LoaderParameters{
				Current: i,
				Total:   t,
				Msg: fmt.Sprintf(
					"Pushing %s to %s...", m.Name,
					path.Join(m.Addr, m.Owner),
				),
				Output: pp.Stdout,
			}),
		},
	)
	if err != nil {
		return err
	}

	req.Header.Add("filename", m.FileName)
	req.Header.Add("email", email)
	req.Header.Add("distro", pp.Distro)
	req.Header.Add("pkgsign", hex.EncodeToString(pkgsign))

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
