// Use of this code is governed by GNU General Public License.
// Official web page: https://ion.lc/core/tab
// Contact email: help@ion.lc

package tab

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"ion.lc/core/tab/msgs"
	"ion.lc/core/tab/pacman"
	"ion.lc/core/tab/process"
)

// Parameters that can be used to build packages.
type BuildParameters struct {
	// Directory where resulting package and signature will be moved.
	Dir string `short:"d" long:"dir" default:"/var/cache/pacman/pkg"`
	// Do not ask for any confirmation on build/installation.
	Quick bool `short:"q" long:"quick"`
	// Syncronize/reinstall package after build.
	Syncbuild bool `short:"s" long:"syncbuild"`
	// Remove dependencies after successful build.
	Rmdeps bool `short:"r" long:"rmdeps"`
	// Do not clean workspace before and after build.
	Dirty bool `short:"g" long:"dirty"`
}

var BuildHelp = `Build, sign and cache package with signature

options:
	-q, --quick     Do not ask for any confirmation (noconfirm)
	-d, --dir <dir> Use custom dir to store result (default /var/cache/pacman/pkg)
	-s, --syncbuild Syncronize dependencies and build target
	-r, --rmdeps    Remove installed dependencies after a successful build
	-g, --dirty     Do not clean workspace before and after build
	-a, --aur       Build targets from AUR git repositories (aur.archlinux.org)

usage:  pack {-B --build} [options] <git/repository(s)>`

// Build package in current directory with provided arguements.
func Build(args []string, prms ...BuildParameters) error {
	p := getOptions(prms)

	msgs.Amsg(os.Stdout, "Building packages")

	msgs.Smsg(os.Stdout, "Running GnuPG check", 1, 2)
	err := CheckGnuPG()
	if err != nil {
		return err
	}

	msgs.Smsg(os.Stdout, "Validating packager identity", 2, 2)
	err = ValidatePackager()
	if err != nil {
		return err
	}

	var builddirs []string
	var buildcurrdir bool

	if len(args) == 0 {
		currdir, err := os.Getwd()
		if err != nil {
			return err
		}
		builddirs = append(args, currdir)
		buildcurrdir = true
	}

	for _, arg := range args {
		dir, err := CloneOrPullDir(os.Stdout, os.Stderr, arg)
		if err != nil {
			return err
		}
		builddirs = append(builddirs, dir)
	}

	for _, dir := range builddirs {
		msgs.Amsg(os.Stdout, "Building package with makepkg")
		err = pacman.Makepkg(pacman.MakepkgParameters{
			Sign:       true,
			Dir:        dir,
			Stdout:     os.Stdout,
			Stderr:     os.Stderr,
			Stdin:      os.Stdin,
			Clean:      !p.Dirty,
			CleanBuild: !p.Dirty,
			Force:      !p.Dirty,
			Install:    p.Syncbuild,
			RmDeps:     p.Rmdeps,
			SyncDeps:   p.Syncbuild,
			Needed:     !p.Syncbuild,
			NoConfirm:  p.Quick,
		})
		if err != nil {
			return errors.Join(err)
		}

		msgs.Amsg(os.Stdout, "Moving package to cache")
		err = CachePackage(dir, p.Dir)
		if err != nil {
			return err
		}
	}

	for _, dir := range builddirs {
		if !buildcurrdir {
			err = os.RemoveAll(dir)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Ensure, that user have created gnupg keys for package signing before package
// is built and cached.
func CheckGnuPG() error {
	hd, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	gpgdir, err := os.ReadDir(path.Join(hd, ".gnupg"))
	if err != nil {
		return errors.New(msgs.ErrGnuPGprivkeyNotFound)
	}
	for _, de := range gpgdir {
		if strings.Contains(de.Name(), "private-keys") {
			return nil
		}
	}
	return errors.New(msgs.ErrGnuPGprivkeyNotFound)
}

// Validate, that packager defined in /etc/makepkg.conf matches signer
// authority in GnuPG.
func ValidatePackager() error {
	keySigner, err := GnuPGidentity()
	if err != nil {
		return err
	}
	f, err := os.ReadFile("/etc/makepkg.conf")
	if err != nil {
		return err
	}
	splt := strings.Split(string(f), "\nPACKAGER=\"")
	if len(splt) != 2 {
		return fmt.Errorf(msgs.ErrNoPackager, keySigner)
	}
	confPackager := strings.Split(splt[1], "\"\n")[0]
	if confPackager != keySigner {
		return fmt.Errorf(msgs.ErrPackagerMissmatch, keySigner, confPackager)
	}
	return nil
}

// Returns name and email from GnuPG. Error, if did not succeed.
func GnuPGidentity() (string, error) {
	cmd := exec.Command("gpg", "-K")
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	err := cmd.Run()
	if err != nil {
		return ``, errors.New("unable to get gnupg identity: " + b.String())
	}
	splt := strings.Split(b.String(), "[ultimate] ")
	if len(splt) < 2 {
		return ``, errors.New("insufficient gpg -K output, unable to get uid")
	}
	return strings.Split(splt[1], "\n")[0], nil
}

// Move package and signature files to cache location defined by user.
func CachePackage(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, de := range entries {
		if strings.HasSuffix(de.Name(), ".pkg.tar.zst") ||
			strings.HasSuffix(de.Name(), ".pkg.tar.zst.sig") {
			err = call(process.Command(&process.Params{
				Sudo:    true,
				Command: "mv",
				Args:    []string{path.Join(src, de.Name()), dst},
			}))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Eject last name from directory or link.
func EjectLastPathArg(s string) string {
	splt := strings.Split(s, "/")
	return splt[len(splt)-1]
}

// This function will clone provided repository to cache directory and return
// name of that directory.
func CloneOrPullDir(outw, errw io.Writer, repo string) (string, error) {
	td := os.TempDir()
	err := os.MkdirAll(path.Join(td, "tab"), os.ModePerm)
	if err != nil {
		return ``, err
	}
	project := EjectLastPathArg(repo)
	msgs.Amsg(outw, "Cloning repository: "+project)
	gitdir := path.Join(td, "tab", project)

	var errbuf bytes.Buffer
	cmd := exec.Command("git", "clone", "https://"+repo, gitdir)
	cmd.Stderr = io.MultiWriter(errw, &errbuf)
	cmd.Stdout = outw
	err = cmd.Run()
	if err != nil {
		if strings.Contains(errbuf.String(), "and is not an empty directory") {
			msgs.Amsg(outw, "Pulling changes")
			err = os.Chdir(gitdir)
			if err != nil {
				return ``, err
			}
			cmd := exec.Command("git", "pull")
			cmd.Stderr = errw
			cmd.Stdout = outw
			return gitdir, cmd.Run()
		}
		return ``, err
	}
	return gitdir, nil
}
