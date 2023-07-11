// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package msgs

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var Help = `Simplified version of pacman

operations:
	pack {-S --sync}   [options] [(registry)/(owner)/package(s)]
	pack {-P --push}   [options] [(registry)/(owner)/package(s)]
	pack {-R --remove} [options] [(registry)/(owner)/package(s)]
	pack {-Q --query}  [options] [(registry)/(owner)/package(s)]
	pack {-B --build}  [options] [git/repository(s)]
	pack {-A --assist} [options] [args]

use 'pack {-h --help}' with an operation for available options`

var SyncHelp = `Syncronize packages

options:
	-q, --quick       Do not ask for any confirmation (noconfirm shortcut)
	-y, --refresh     Download fresh package databases from the server (-yy force)
	-u, --upgrade     Upgrade installed packages (-uu enables downgrade)
	-f, --force       Reinstall up to date targets

usage:  pack {-S --sync} [options] <(registry)/(owner)/package(s)>`

var PushHelp = `Push packages

options:
	-d, --dir <dir> Use custom source dir with packages (default pacman cache)
	-w, --insecure  Push package over HTTP instead of HTTPS
	    --distro    Assign custom distribution in registry (default archlinux)
	    --endpoint  Use custom API endpoints rootpath

usage:  pack {-P --push} [options] <registry/(owner)/package(s)>`

var RemoveHelp = `Remove packages

options:
	-c, --confirm  Ask for confirmation when deleting package
	-a, --norecurs Leave package dependencies in the system (removed by default)
	-j, --nocfgs   Leave package configs in the system (removed by default)
	-k, --cascade  Remove packages and all packages that depend on them
	    --arch     Custom architecture for remote delete operation
		--distro   Custom distribution for remote delete operation

usage:  pack {-R --remove} [options] <(registry)/(owner)/package(s)>`

var QueryHelp = `Query packages

options:
	-i, --info     View package information (-ii for backup files)
	-l, --list     List the files owned by the queried package
	-o, --outdated List outdated packages

usage:  pack {-Q --query} [options] <(registry)/(owner)/package(s)>`

var BuildHelp = `Build package

options:
	-q, --quick     Do not ask for any confirmation (noconfirm)
	-d, --dir <dir> Use custom dir to store result (default /var/cache/pacman/pkg)
	-s, --syncbuild Syncronize dependencies and build target
	-r, --rmdeps    Remove installed dependencies after a successful build
	-g, --garbage   Do not clean workspace before and after build

usage:  pack {-B --build} [options] <git/repository(s)>`

var AssistHelp = `Additional utilities

options:
	-e, --export  Export public GnuPG key armor
	-n, --gen     Generate GnuPG key for package singing
	    --recv    Run key recieve operaion
	    --info    Show information about your GnuPG keys
	    --setpkgr Automatically set packager in makepkg.conf
	    --flutter Generate PKGBUILD, app.sh and app.desktop for flutter application
	    --gocli   Generate PKGBUILD for CLI utility in go

usage:  pack {-A --assist} [options] <(args)>`

var Version = `             Pack - package manager.
          Copyright (C) 2023 FMNX team
     
  This program may be freely redistributed under
   the terms of the GNU General Public License.
       Web page: https://fmnx.su/core/pack
 
                 Version: 0.1.1`

var Color bool

func init() {
	b, err := os.ReadFile("/etc/pacman.conf")
	if err != nil {
		fmt.Println("unable to read pacman configuration")
		os.Exit(1)
	}
	Color = strings.Contains(string(b), "\nColor\n")
	if !Color {
		color.NoColor = true
	}
}
