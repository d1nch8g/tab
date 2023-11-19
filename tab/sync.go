// Use of this code is governed by GNU General Public License.
// Official web page: https://ion.lc/core/tab
// Contact email: help@ion.lc

package tab

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"ion.lc/core/tab/pacman"
	"ion.lc/core/tab/process"
)

type SyncParameters struct {
	// Download fresh package databases from the server (-yy force)
	Refresh []bool `short:"y" long:"refresh"`
	// Upgrade installed packages (-uu enables downgrade)
	Upgrade []bool `short:"u" long:"upgrade"`
	// Don't ask for any confirmation (--noconfirm)
	Quick bool `short:"q" long:"quick"`
	// Reinstall up to date targets
	Force bool `short:"f" long:"force"`
	// Use HTTP instead of https
	Insecure bool `short:"i" long:"insecure"`
}

var SyncHelp = `Syncronize packages

options:
	-q, --quick    Do not ask for any confirmation (noconfirm shortcut)
	-y, --refresh  Download fresh package databases from the server (-yy force)
	-u, --upgrade  Upgrade installed packages (-uu enables downgrade)
	-f, --force    Reinstall up to date targets
	-i, --insecure Use HTTP protocol for new pacman databases (HTTPS by default)

usage: tab {-S --sync} [options] <(registry)/(owner)/package(s)>`

// Syncronize provided packages with provided parameters.
func Sync(args []string, prms ...SyncParameters) error {
	p := getParameters(prms)

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
			Stdout:    os.Stdout,
			Stderr:    os.Stderr,
			Stdin:     os.Stdin,
		})
	}

	conf, err = addMissingDatabases(args, p.Insecure)
	if err != nil {
		return err
	}

	pkgs = formatPackages(args)

	err = pacman.SyncList(pkgs, pacman.SyncParameters{
		Sudo:      true,
		Needed:    !p.Force,
		NoConfirm: p.Quick,
		Refresh:   p.Refresh,
		Upgrade:   p.Upgrade,
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		Stdin:     os.Stdin,
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
			if strings.Contains(conf, fmt.Sprintf("[%s]", splt[0])) {
				continue
			}
			addConfDatabase(protocol, splt[0], splt[0], "")
		case 3:
			if strings.Contains(conf, fmt.Sprintf("[%s.%s]", splt[1], splt[0])) {
				continue
			}
			addConfDatabase(protocol, splt[1]+"."+splt[0], splt[0], "/"+splt[1])
		}
	}
	return &conf, nil
}

// Simple function to add database to pacman.conf.
func addConfDatabase(protocol, database, domain, owner string) error {
	const (
		confroot = "\n[%s]\nServer = %s://%s/api/packages/%s/arch/%s/%s\n"
		os       = "archlinux"
		arch     = "x86_64"
	)
	tmpl := fmt.Sprintf(confroot, database, protocol, domain, owner, os, arch)
	command := "cat <<EOF >> /etc/pacman.conf" + tmpl + "EOF"
	return call(process.Command(&process.Params{
		Sudo:    true,
		Command: "bash",
		Args:    []string{"-c", command},
	}))
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
	return call(process.Command(&process.Params{
		Sudo:    true,
		Command: "bash",
		Args:    []string{"-c", "cat <<EOF > /etc/pacman.conf\n" + s + "EOF"},
	}))
}
