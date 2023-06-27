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
	"fmnx.su/core/pack/pacman"
)

type SyncParameters struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	// Download fresh package databases from the server (-yy force)
	Refresh []bool
	// Upgrade installed packages (-uu enables downgrade)
	Upgrade []bool
	// Don't ask for any confirmation (--noconfirm)
	Quick bool
	// Reinstall up to date targets
	Force bool
	// Use HTTP instead of https
	Insecure bool
}

func syncdefault() *SyncParameters {
	return &SyncParameters{
		Quick:   true,
		Refresh: []bool{true},
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Stdin:   os.Stdin,
	}
}

// Syncronize provided packages with provided parameters.
func Sync(args []string, prms ...SyncParameters) error {
	p := formOptions(prms, syncdefault)

	var err error
	var conf *string
	var pkgs []string

	if len(args) == 0 {
		return pacman.SyncList(pkgs, pacman.SyncParameters{
			Sudo:      true,
			Needed:    !p.Force,
			NoConfirm: p.Quick,
			Refresh:   p.Refresh,
			Upgrade:   p.Upgrade,
			Stdout:    p.Stdout,
			Stderr:    p.Stderr,
			Stdin:     p.Stdin,
		})
	}

	msgs.Amsg(p.Stdout, "Syncronizing packages")

	msgs.Smsg(p.Stdout, "Adding missing databases to pacman.conf", 1, 2)
	conf, err = addMissingDatabases(args, p.Insecure)
	if err != nil {
		return err
	}

	msgs.Smsg(p.Stdout, "Preparing packages to sync format", 2, 2)
	pkgs = formatPackages(args)

	err = pacman.SyncList(pkgs, pacman.SyncParameters{
		Sudo:      true,
		Needed:    !p.Force,
		NoConfirm: p.Quick,
		Refresh:   p.Refresh,
		Upgrade:   p.Upgrade,
		Stdout:    p.Stdout,
		Stderr:    p.Stderr,
		Stdin:     p.Stdin,
	})
	if err != nil {
		return errors.Join(err, writeconf(*conf))
	}
	return nil
}

// Iterate over packages, check wether package database is present, if not
// add new database to pacman.conf. Return previous version of pacman.conf.
func addMissingDatabases(pkgs []string, insecure bool) (*string, error) {
	protocol := "https"
	if insecure {
		protocol = "http"
	}
	f, err := os.ReadFile("/etc/pacman.conf")
	if err != nil {
		return nil, err
	}
	conf := string(f)
	for _, pkg := range pkgs {
		splt := strings.Split(pkg, "/")
		switch len(splt) {
		case 2:
			if strings.Contains(conf, splt[0]+"/api/packages/arch") {
				continue
			}
			addConfDatabase(protocol, splt[0], splt[0], "")
		case 3:
			if strings.Contains(conf, splt[0]+"/api/packages/arch/"+splt[1]) {
				continue
			}
			addConfDatabase(protocol, splt[1]+"."+splt[0], splt[0], "/"+splt[1])
		}
	}
	return &conf, nil
}

// Simple function to add database to pacman.conf.
func addConfDatabase(protocol, database, domain, owner string) error {
	const confroot = "\n[%s]\nServer = %s://%s/api/packages/%s/arch/%s/%s\n"
	os := "archlinux"
	tmpl := fmt.Sprintf(confroot, database, protocol, domain, owner, os, "x86_64")
	command := "cat <<EOF >> /etc/pacman.conf" + tmpl + "EOF"
	return call(exec.Command("sudo", "bash", "-c", command))
}

// Format packages to pre-sync format.
func formatPackages(pkgs []string) []string {
	var out []string
	for _, pkg := range pkgs {
		splt := strings.Split(pkg, "/")
		switch len(splt) {
		case 1:
			out = append(out, pkg)
		case 2:
			out = append(out, splt[0]+"/"+splt[1])
		case 3:
			out = append(out, splt[1]+"."+splt[0]+"/"+splt[2])
		}
	}
	return out
}

// Overwrite pacman.conf with provided string.
func writeconf(s string) error {
	return call(exec.Command(
		"sudo", "bash", "-c",
		"cat <<EOF > /etc/pacman.conf\n"+s+"EOF",
	))
}
