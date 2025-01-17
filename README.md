<p align="center">
<img style="align: center; padding-left: 10px; padding-right: 10px; padding-bottom: 10px;" width="238px" height="238px" src="./logo.png" />
</p>

<h2 align="center">Tab - decentralized package manager</h2>

![](https://img.shields.io/badge/alpha-0.2.0-blue.svg)
[![](https://img.shields.io/badge/license-GPL-orange.svg)](https://www.gnu.org/licenses/gpl-3.0.en.html)
[![](https://img.shields.io/badge/arch-package-00bcd4.svg)](https://ion.lc/core/-/packages/arch/tab)

Tab is decentralized package manager - automates the process of installing, upgrading, configuring, and removing programs, not tying user to specific package registry. Simple web links `example.com/owner/package` can be used to install software.

Tab works as a wrapper over [pacman](https://wiki.archlinux.org/title/Pacman), keeping full compatability with [arch package format](https://wiki.archlinux.org/title/Arch_package_guidelines) (typically, `.pkg.tar.zst` packages with [GnuPG](https://gnupg.org/) signatures). Also tab provides additional functionality for software delivery and handles pacman database management automatically.

Main goal of the project is to simplify process of arch package creation, increase delivery speed and create simple and accesible but still powerfull package management solution.

---

### Installation

Single line installation script for all arch-based distributions:

```sh
git clone https://ion.lc/core/tab && cd tab && makepkg -sfri --needed --noconfirm
```

---

### Operations

1. Sync packages - operation that you use to install packages to the system.

```sh
tab -Sy nano blender example.com/owner/package
```

You can mix packages with and without registries in input. This command will add missing registries to `pacman.conf` and try to synchronize packages with pacman. Flags for operation:

- `-q`, `--quick` - Do not ask for any confirmation (noconfirm shortcut)
- `-y`, `--refresh` - Download fresh package databases from the server (-yy force)
- `-u`, `--upgrade` - Upgrade installed packages (-uu enables downgrade)
- `-f`, `--force` - Reinstall up to date targets
- `-i`, `--insecure` - Use HTTP protocol for new pacman databases (HTTPS by default)

2. Query packages - operation that you use to inspect the state of your system or view package parameters.

```sh
tab -Qi pacman
```

- `-i`, `--info` - View package information (-ii for backup files)
- `-l`, `--list` - List the files owned by the queried package
- `-o`, `--outdated` - List outdated packages

3. Remove packages - this operation will remove packages from the system or registry. By default, it removes local packages, if you provide a registry, remote deletion will be executed. When removing remote packages, they provide a version after @.

```sh
tab -R vim
tab -R example.com/owner/package@1-1
```

- `-c`, `--confirm` - Ask for confirmation when deleting package
- `-r`, `--norecurs` - Leave package dependencies in the system (removed by default)
- `-f`, `--nocfgs` - Leave package configs in the system (removed by default)
- `-c`, `--cascade` - Remove packages and all packages that depend on them
- `-i`, `--insecure` - Use HTTP protocol for API calls (remote delete)

4. Build packages - command that you use to build packages. If you provide git repo(s) in arguments, this command will clone and build them.

```sh
tab -B aur.archlinux.org/veloren-bin ion.lc/core/ainst
tab -Bqsa onlyoffice-bin
```

After a successful build, prepared packages are stored in `/var/cache/pacman/pkg`. Delete flags:

- `-q`, `--quick` - Do not ask for any confirmation (noconfirm)
- `-d`, `--dir` - Use custom dir to store result (default /var/cache/pacman/pkg)
- `-s`, `--syncbuild` - Syncronize dependencies and build target
- `-r`, `--rmdeps` - Remove installed dependencies after a successful build
- `-g`, `--dirty` - Do not clean workspace before and after build
- `-a`, `--aur` - Build targets from AUR git repositories (aur.archlinux.org)

5. Push packages - operation that you use to deliver your software to any pack registry (currently only gitea supported).

```sh
tab -P ion.lc/core/onlyoffice-bin
```

- `-d`, `--dir` - Use custom source dir with packages (default pacman cache)
- `-i`, `--insecure` - Push package over HTTP instead of HTTPS
- `-s`, `--distro` - Assign custom distribution in registry (default archlinux)
- `-e`, `--export` - Export public GPG key armor
