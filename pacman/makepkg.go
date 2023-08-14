// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pacman

import (
	"io"
	"os"

	"fmnx.su/core/pack/process"
)

// Options for building packages.
type MakepkgParameters struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	// Additional parameters, that will be appended to command as arguements.
	AdditionalParams []string
	// Directory where process will be executed.
	Dir string
	// Specify a key to use for gpg signing instead of the default. [--key <key>]
	GpgKey string
	// Use an alternate build script (not 'PKGBUILD'). [-p <file>]
	File string
	// Ignore incomplete arch field in PKGBUILD. [--ignorearch]
	IgnoreEach bool
	// Clean up work files after build. [--clean]
	Clean bool
	// Remove $srcdir/ dir before building the package. [--cleanbuild]
	CleanBuild bool
	// Skip all dependency checks. [--nodeps]
	NoDeps bool
	// Do not extract source files (use existing $srcdir/ dir). [--noextract]
	NoExtract bool
	// Overwrite existing package. [--force]
	Force bool
	// Generate integrity checks for source files. [--geninteg]
	Geinteg bool
	// Install package after successful build. [--install]
	Install bool
	// Log package build process. [--log]
	Log bool
	// Disable colorized output messages. [--nocolor]
	NoColor bool
	// Download and extract files only. [--nobuild]
	NpBuild bool
	// Remove installed dependencies after a successful build. [--rmdeps]
	RmDeps bool
	// Repackage contents of the package without rebuilding. [--repackage]
	Repackage bool
	// Install missing dependencies with pacman. [--syncdeps]
	SyncDeps bool
	// Use an alternate config file (not '/etc/makepkg.conf'). [--config <file>]
	Config string
	// Do not update VCS sources. [--holdver]
	HoldVer bool
	// Do not create package archive. [--noarchive]
	NoArchive bool
	// Do not run the check() function in the PKGBUILD. [--nocheck]
	NoCheck bool
	// Do not run the prepare() function in the PKGBUILD. [--noprepare]
	NoPrepare bool
	// Do not create a signature for the package. [--nosign]
	NoSign bool
	// Sign the resulting package with gpg. [--sign]
	Sign bool
	// Do not verify checksums of the source files. [--skipchecksums]
	SkipCheckSums bool
	// Do not perform any verification checks on source files. [--skipinteg]
	SkipIntegrityChecks bool
	// Do not verify source files with GPG signatures. [--skippgpcheck]
	SkipPgpCheck bool
	// Do not reinstall up to date packages. [--needed]
	Needed bool
	// Do not ask for any confirmation. [--noconfirm]
	NoConfirm bool
	// Do not show a progress bar when downloading files. [--noprogressbar]
	NoProgressBar bool
	// Install packages as non-explicitly installed. [--asdeps]
	AsDeps bool
}

func makepkgdefault() *MakepkgParameters {
	return &MakepkgParameters{
		Clean:      true,
		Force:      true,
		Sign:       true,
		CleanBuild: true,
		Install:    true,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		Stdin:      os.Stdin,
	}
}

// This command will build a package in directory provided in options.
// Function is safe for concurrent usage. Can be called from multiple
// goruotines, when options Install or SyncDeps are false.
func Makepkg(opts ...MakepkgParameters) error {
	p := formOptions(opts, makepkgdefault)

	var args []string
	if p.IgnoreEach {
		args = append(args, "--ignorearch")
	}
	if p.Clean {
		args = append(args, "--clean")
	}
	if p.CleanBuild {
		args = append(args, "--cleanbuild")
	}
	if p.NoDeps {
		args = append(args, "--nodeps")
	}
	if p.NoExtract {
		args = append(args, "--noextract")
	}
	if p.Force {
		args = append(args, "--force")
	}
	if p.Geinteg {
		args = append(args, "--geninteg")
	}
	if p.Log {
		args = append(args, "--log")
	}
	if p.NoColor {
		args = append(args, "--nocolor")
	}
	if p.NpBuild {
		args = append(args, "--nobuild")
	}
	if p.RmDeps {
		args = append(args, "--rmdeps")
	}
	if p.Repackage {
		args = append(args, "--repackage")
	}
	if p.HoldVer {
		args = append(args, "--holdver")
	}
	if p.NoArchive {
		args = append(args, "--noarchive")
	}
	if p.NoCheck {
		args = append(args, "--nocheck")
	}
	if p.NoPrepare {
		args = append(args, "--noprepare")
	}
	if p.NoSign {
		args = append(args, "--nosign")
	}
	if p.Sign {
		args = append(args, "--sign")
	}
	if p.SkipCheckSums {
		args = append(args, "--skipchecksums")
	}
	if p.SkipIntegrityChecks {
		args = append(args, "--skipinteg")
	}
	if p.SkipPgpCheck {
		args = append(args, "--skippgpcheck")
	}
	if p.Needed {
		args = append(args, "--needed")
	}
	if p.NoConfirm {
		args = append(args, "--noconfirm")
	}
	if p.NoProgressBar {
		args = append(args, "--noprogressbar")
	}
	if p.AsDeps {
		args = append(args, "--asdeps")
	}
	if p.File != `` {
		args = append(args, "-p")
		args = append(args, p.File)
	}
	if p.Config != "" {
		args = append(args, "--config")
		args = append(args, p.Config)
	}
	if p.GpgKey != "" {
		args = append(args, "--key")
		args = append(args, p.GpgKey)
	}
	if p.Install {
		args = append(args, "--install")
		mu.Lock()
		defer mu.Unlock()
	}
	if p.SyncDeps {
		args = append(args, "--syncdeps")
		if mu.TryLock() {
			defer mu.Unlock()
		}
	}
	args = append(args, p.AdditionalParams...)

	return process.Command(&process.Params{
		Stdout:  p.Stdout,
		Stderr:  p.Stderr,
		Stdin:   p.Stdin,
		Command: makepkg,
		Args:    args,
		Dir:     p.Dir,
	}).Run()
}
