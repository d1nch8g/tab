<p align="center">
<img style="align: center; padding-left: 10px; padding-right: 10px; padding-bottom: 10px;" width="238px" height="238px" src="./logo.png" />
</p>

<h2 align="center">Pack - package manager</h2>

[![Generic badge](https://img.shields.io/badge/license-gpl-orange.svg)](https://fmnx.su/core/pack/src/branch/main/LICENSE)
[![Generic badge](https://img.shields.io/badge/fmnx-repo-006db0.svg)](https://fmnx.su/core/pack)
[![Generic badge](https://img.shields.io/badge/codeberg-repo-45a3fb.svg)](https://codeberg.org/fmnx/pack)
[![Generic badge](https://img.shields.io/badge/github-repo-white.svg)](https://github.com/fmnx-io/pack)
[![Generic badge](https://img.shields.io/badge/package-0.1.0_alpha-00bcd4.svg)](https://fmnx.su/core/-/packages/arch/pack)

Pack works as a wrapper over pacman providing additional functionality for software delivery and connection management.

For users pack provides ability to install packages from any compatible registry using package URL. For developers pack is offering simple interface for quick package delivery.

Share and install software without any intermediate regulator.

---

### Installation

Single line installation script for all arch based distributions:

```sh
git clone https://fmnx.su/core/pack && cd pack && makepkg -sfri --needed --noconfirm
```

Alternatively, you can install pack with `go`:

```sh
go install fmnx.su/core/pack
```

---

### Operations

1. Sync packages - operation that you use to install packages to the system. You can mix packages with and without registries in command input. This command will add missing registries to `pacman.conf` and try to syncronize packages with pacman.

```sh
⚡ Syncronize packages

options:
	-q, --quick       Do not ask for any confirmation (noconfirm shortcut)
	-y, --refresh     Download fresh package databases from the server (-yy force)
	-u, --upgrade     Upgrade installed packages (-uu enables downgrade)
	-f, --force       Reinstall up to date targets

usage:  pack {-S --sync} [options] <(registry)/(owner)/package(s)>
```

2. Push packages - operation that you use to deliver your software to any pack registry (standalone registry or gitea). Registry parameter is required, owner paarameter is optional.

```sh
🚀 Push packages

options:
        -d, --dir <dir> Use custom source dir with packages (default pacman cache)
        -w, --insecure  Push package over HTTP instead of HTTPS
            --distro    Assign custom distribution in registry (default archlinux)
            --endpoint  Use custom API endpoints rootpath

usage:  pack {-P --push} [options] <registry/(owner)/package(s)>
```

3. Remove packages - this operation will remove packages from system or remote depending on provided arguement. If reigsty and owner are provided, then remote deletion will be executed, otherwise package will be deleted on local system.

```sh
📍 Remove packages

options:
        -o, --confirm  Ask for confirmation when deleting package
        -a, --norecurs Leave package dependencies in the system (removed by default)
        -w, --nocfgs   Leave package configs in the system (removed by default)
            --cascade  Remove packages and all packages that depend on them

usage:  pack {-R --remove} [options] <package(s)>
```

4. Query packages - this command can be executed to get information about local or remote packages. For targets without registry and owner specified local description will be provided, for targets with registry remote information.
<!-- If you want to search for a package on remote, just put @ before target package -->

```sh
🔎 Query packages

options:
        -i, --info     View package information (-ii for backup files)
        -l, --list     List the files owned by the queried package
        -o, --outdated List outdated packages

usage:  pack {-Q --query} [options] <(registry)/(owner)/package(s)>
```

5. Build packages - command that will build package in current directory if no arguements provided, otherwise it will treat packages as git repositories, clone them to `~/.packcache` directory, build and remove directory afterwards.

```sh
🔐 Build package

options:
        -q, --quick     Do not ask for any confirmation (noconfirm)
        -d, --dir <dir> Use custom dir to store result (default /var/cache/pacman/pkg)
        -s, --syncbuild Syncronize dependencies and build target
        -r, --rmdeps    Remove installed dependencies after a successful build
        -g, --garbage   Do not clean workspace before and after build

usage:  pack {-B --build} [options] <(registry)/(owner)/package(s)>
```

6. Util - small assistant that can generate generate tempaltes or help with GnuPG related operations.

```sh
📄 Additional utilities

options:
        --info    Show information about your GnuPG keys
        --gen     Generate GnuPG key for package singing
        --export   Export public GnuPG key armor
        --recv    Run recieve key operaion
        --setpkgr Automatically set packager in makepkg.conf
        --flutter Generate PKGBUILD, app.sh and app.desktop for flutter application
        --gocli   Generate PKGBUILD for CLI utility in go

usage:  pack {-U --util} [options] <(args)>
```
