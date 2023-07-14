<p align="center">
<img style="align: center; padding-left: 10px; padding-right: 10px; padding-bottom: 10px;" width="238px" height="238px" src="./logo.png" />
</p>

<h2 align="center">Pack - package manager</h2>

![Generic badge](https://img.shields.io/badge/status-alpha-red.svg)
[![Generic badge](https://img.shields.io/badge/license-gpl-orange.svg)](https://fmnx.su/core/pack/src/branch/main/LICENSE)
[![Generic badge](https://img.shields.io/badge/fmnx-repo-006db0.svg)](https://fmnx.su/core/pack)
[![Generic badge](https://img.shields.io/badge/codeberg-repo-45a3fb.svg)](https://codeberg.org/fmnx/pack)
[![Generic badge](https://img.shields.io/badge/github-repo-white.svg)](https://github.com/fmnx-su/pack)
[![Generic badge](https://img.shields.io/badge/arch-package-00bcd4.svg)](https://fmnx.su/core/-/packages/arch/pack)

> **Warning!** Project is in alpha stage, API's are very likely to be changed.

Pack works as a wrapper over pacman providing additional functionality for software delivery and pacman database management.

For users pack provides ability to install packages from any compatible registry using package URL. For developers pack is offering simple interface for quick package delivery.

---

### Installation

Single line installation script for all arch based distributions:

```sh
git clone https://fmnx.su/core/pack && cd pack && makepkg -sfri --needed --noconfirm
```

---

### Operations

1. Sync packages - operation that you use to install packages to the system.

```sh
pack -Sy nano blender example.com/owner/package
```

You can mix packages with and without registries in input. This command will add missing registries to `pacman.conf` and try to syncronize packages with pacman. Flags for operation:

- `-q`, `--quick` - Do not ask for any confirmation (noconfirm shortcut)
- `-y`, `--refresh` - Download fresh package databases from the server (-yy force)
- `-u`, `--upgrade` - Upgrade installed packages (-uu enables downgrade)
- `-f`, `--force` - Reinstall up to date targets

2. Query packages - operation that you use to inspect the state of your system or view package parameters.

```sh
pack -Qi pacman
```

- `-i`, `--info` - View package information (-ii for backup files)
- `-l`, `--list` - List the files owned by the queried package
- `-o`, `--outdated` - List outdated packages

3. Remove packages - this operation will remove packages from system or registry. By default removes local packages, if you provide registry remote deletion will be executed. When removing remote packages provide version after @.

```sh
pack -R vim
pack -R example.com/owner/package@1-1
```

- `-o`, `--confirm` - Ask for confirmation when deleting package
- `-a`, `--norecurs` - Leave package dependencies in the system (removed by default)
- `-w`, `--nocfgs` - Leave package configs in the system (removed by default)
- `-k`, `--cascade` - Remove packages and all packages that depend on them
- `--arch` - Custom architecture for remote delete operation
- `--distro` - Custom distribution for remote delete operation

4. Build packages - command that you use to build packages. If you provide git repo(s) in args, this command will clone and build them.

```sh
pack -B aur.archlinux.org/veloren-bin aur.archlinux.org/onlyoffice-bin
```

After successful build prepared packages are stored in `/var/cache/pacman/pkg`. Delete flags:

- `-q`, `--quick` - Do not ask for any confirmation (noconfirm)
- `-d`, `--dir` _directory_ - Use custom dir to cache built package (default /var/cache/pacman/pkg)
- `-s`, `--syncbuild` - Syncronize dependencies and build target
- `-r`, `--rmdeps` - Remove installed dependencies after a successful build
- `-g`, `--garbage` - Do not clean workspace before and after build

5.  Push packages - operation that you use to deliver your software to any pack registry (currently standalone registry or gitea).

```sh
pack -P fmnx.su/core/onlyoffice-bin
```

- `-w`, `--insecure` - Push package over HTTP instead of HTTPS
- `-d`, `--dir` _directory_ - Use custom source dir with packages (default pacman cache)
- `--distro` - Assign custom distribution in registry (default archlinux)
- `--endpoint` - Use custom API endpoints rootpath

6. Assist - generate project tempaltes, export [GnuPG](https://gnupg.org/) keys, set packages, etc...

```sh
pack -A
```

- `-e`, `--export` - Export public GnuPG key armor
- `-x`, `--fix` - Check/fix compatability of identities in git, gpg and makepkg.
- `--flutter` - Generate PKGBUILD, app.sh and app.desktop for flutter application
- `--gocli` - Generate PKGBUILD for CLI utility in go
