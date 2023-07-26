<p align="center">
<img style="align: center; padding-left: 10px; padding-right: 10px; padding-bottom: 10px;" width="238px" height="238px" src="./logo.png" />
</p>

<h2 align="center">Pack - package manager</h2>

![](https://img.shields.io/badge/status-alpha-red.svg)
[![](https://img.shields.io/badge/license-gpl-orange.svg)](https://fmnx.su/core/pack/src/branch/main/LICENSE)
[![](https://img.shields.io/badge/fmnx-repo-006db0.svg)](https://fmnx.su/core/pack)
[![](https://img.shields.io/badge/codeberg-repo-45a3fb.svg)](https://codeberg.org/fmnx/pack)
[![](https://img.shields.io/badge/github-repo-white.svg)](https://github.com/fmnx-su/pack)
[![](https://img.shields.io/badge/arch-package-00bcd4.svg)](https://fmnx.su/core/-/packages/arch/pack)

> **Warning!** The project is in the alpha stage, so APIs are very likely to be changed.

Pack works as a wrapper over pacman, providing additional functionality for software delivery and pacman database management.

For users, pack provides the ability to install packages from any compatible registry using package URL. For developers, pack the offers a simple interface for quick package delivery.

---

### Installation

Single line installation script for all arch-based distributions:

```sh
git clone https://fmnx.su/core/pack && cd pack && makepkg -sfri --needed --noconfirm
```

---

### Operations

1. Sync packages - operation that you use to install packages to the system.

```sh
pack -Sy nano blender for example.com/owner/package
```

You can mix packages with and without registries in input. This command will add missing registries to `pacman.conf` and try to synchronize packages with pacman. Flags for operation:

- `-q`, `--quick` - Do not ask for any confirmation (non-confirm shortcut)
- `-y`, `--refresh` - Download fresh package databases from the server (-yy force)
- `-u`, `--upgrade` - Upgrade installed packages (-uu enables downgrade)
- `-f`, `--force` - Reinstall up-to-date targets

2. Query packages - operation that you use to inspect the state of your system or view package parameters.

```sh
pack -Qi pacman
```

- `-i`, `--info` - View package information (-ii for backup files)
- `-l`, `--list` - List the files owned by the queried package
- `-o`, `--outdated` - List outdated packages

3. Remove packages - this operation will remove packages from the system or registry. By default, it removes local packages, if you provide a registry, remote deletion will be executed. When removing remote packages, they provide a version after @.

```sh
pack -R vim
pack -R for example.com/owner/package@1-1
```

- `-o`, `--confirm` - Ask for confirmation when deleting a package
- `-a`, `--norecurs` - Leave package dependencies in the system (removed by default)
- `-w`, `--nocfgs` - Leave package configs in the system (removed by default)
- `-k`, `--cascade` - Remove packages and all packages that depend on them
- `--arch` - Custom architecture for remote delete operation
- `--distro` - Custom distribution for remote delete operation

4. Build packages - command that you use to build packages. If you provide git repo(s) in arguments, this command will clone and build them.

```sh
pack -B aur.archlinux.org/veloren-bin aur.archlinux.org/onlyoffice-bin
```

After a successful build, prepared packages are stored in `/var/cache/pacman/pkg`. Delete flags:

- `-q`, `--quick` - Do not ask for any confirmation (confirm)
- `-d`, `--dir` _directory_ - Use custom directory to cache built package (default /var/cache/pacman/pkg)
- `-s`, `--syncbuild` - Synchronize dependencies and build target
- `-r`, `--rmdeps` - Remove installed dependencies after a successful build
- `-g`, `--garbage` - Do not clean the workspace before and after building

5. Push packages - operation that you use to deliver your software to any pack registry (currently standalone registry or gitea).

```sh
pack -P fmnx.su/core/onlyoffice-bin
```

- `-w`, `--insecure` - Push package over HTTP instead of HTTPS
- `-d`, `--dir` _directory_ - Use custom source directory with packages (default pacman cache)
- `--distro` - Assign custom distribution in registry (default archlinux)
- `--endpoint` - Use custom API endpoints root path

6. Assist - generate project templates, export [GnuPG](https://gnupg.org/) keys, set packages, etc...

```sh
pack -A
```

- `-e`, `--export` - Export public GnuPG key armor
- `-x`, `--fix` - Check/fix compatibility of identities in git, gpg and makepkg.
- `--flutter` - Generate PKGBUILD, app.sh and app.desktop for flutter application
- `--gocli` - Generate PKGBUILD for CLI utility in go
