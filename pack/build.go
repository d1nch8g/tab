// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"fmnx.su/core/pack/msgs"
	"fmnx.su/core/pack/pacman"
)

// Parameters that can be used to build packages.
type BuildParameters struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	// Directory where resulting package and signature will be moved.
	Dir string
	// Do not ask for any confirmation on build/installation.
	Quick bool
	// Syncronize/reinstall package after build.
	Syncbuild bool
	// Remove dependencies after successful build.
	Rmdeps bool
	// Do not clean workspace before and after build.
	Garbage bool
}

func builddefault() *BuildParameters {
	return &BuildParameters{
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		Stdin:     os.Stdin,
		Dir:       "/var/cache/pacman/pkg",
		Syncbuild: true,
		Rmdeps:    true,
	}
}

// Build package in current directory with provided arguements
func Build(args []string, prms ...BuildParameters) error {
	p := formOptions(prms, builddefault)

	msgs.Amsg(p.Stdout, "Building packages")

	msgs.Smsg(p.Stdout, "Running GnuPG check", 1, 2)
	err := checkGnupg()
	if err != nil {
		return err
	}

	msgs.Smsg(p.Stdout, "Validating packager identity", 2, 2)
	err = validatePackager()
	if err != nil {
		return err
	}

	var builddirs []string

	if len(args) == 0 {
		currdir, err := os.Getwd()
		if err != nil {
			return err
		}
		builddirs = append(args, currdir)
	}

	for _, arg := range args {
		dir, err := cloneOrPullDir(p.Stdout, p.Stderr, arg)
		if err != nil {
			return err
		}
		builddirs = append(builddirs, dir)
	}

	for _, dir := range builddirs {
		msgs.Amsg(p.Stdout, "Building package with makepkg")
		err = pacman.Makepkg(pacman.MakepkgParameters{
			Sign:       true,
			Dir:        dir,
			Stdout:     p.Stdout,
			Stderr:     p.Stderr,
			Stdin:      p.Stdin,
			Clean:      !p.Garbage,
			CleanBuild: !p.Garbage,
			Force:      !p.Garbage,
			Install:    p.Syncbuild,
			RmDeps:     p.Rmdeps,
			SyncDeps:   p.Syncbuild,
			Needed:     !p.Syncbuild,
			NoConfirm:  p.Quick,
		})
		if err != nil {
			return errors.Join(err)
		}

		msgs.Amsg(p.Stdout, "Moving package to cache")
		movecommand := "sudo mv " + dir + "/*.pkg.tar.zst* " + p.Dir
		cmd := exec.Command("bash", "-c", movecommand)
		err = call(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

// Ensure, that user have created gnupg keys for package signing before package
// is built and cached.
func checkGnupg() error {
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
func validatePackager() error {
	keySigner, err := gnuPGIdentity()
	if err != nil {
		return err
	}
	f, err := os.ReadFile("/etc/makepkg.conf")
	if err != nil {
		return err
	}
	splt := strings.Split(string(f), "\nPACKAGER=\"")
	if len(splt) != 2 {
		return errors.New(msgs.ErrNoPackager)
	}
	confPackager := strings.Split(splt[1], "\"\n")[0]
	if confPackager != keySigner {
		return errors.New(msgs.ErrSignerMissmatch)
	}
	return nil
}

// Returns name and email from GnuPG. Error, if did not succeed.
func gnuPGIdentity() (string, error) {
	gnukey := `gpg --with-colons -k | awk -F: '$1=="uid" {print $10; exit}'`
	cmd := exec.Command("bash", "-c", gnukey)
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	err := cmd.Run()
	if err != nil {
		o := b.String()
		return ``, errors.New("unable to get gnupg identity: " + o)
	}
	return strings.ReplaceAll(b.String(), "\n", ""), nil
}

// Eject last name from directory or link.
func ejectLastPathArg(s string) string {
	splt := strings.Split(s, "/")
	return splt[len(splt)-1]
}

// This function will clone provided repository to cache directory and return
// name of that directory.
func cloneOrPullDir(outw, errw io.Writer, repo string) (string, error) {
	uhd, err := os.UserHomeDir()
	if err != nil {
		return ``, err
	}
	err = os.MkdirAll(path.Join(uhd, ".packcache"), os.ModePerm)
	if err != nil {
		return ``, err
	}
	project := ejectLastPathArg(repo)
	msgs.Amsg(outw, "Cloning repository: "+project)
	gitdir := path.Join(uhd, ".packcache", project)

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
