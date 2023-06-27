// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pacman

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"strings"
)

// Query parameters for pacman packages.
type QueryParameters struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	// List packages explicitly installed.
	Explicit bool
	// List packages installed as dependencies.
	Deps bool
	// Query only for packages installed from official repositories.
	Native bool
	// Query for packages installed from other sources.
	Foreign bool
	// Unrequired packages (not a dependency for other one).
	Unrequired bool
	// View all members of a package group.
	Groups bool
	// View package information (-ii for backup files).
	Info []bool
	// Check that package files exist (-kk for file properties).
	Check []bool
	// List the files owned by the queried package.
	List []bool
	// Query a package file instead of the database.
	File string
	// List packages to upgrade. [-u]
	Upgrade bool
	// Additional queue parameters.
	AdditionalParams []string
}

func QueryDefault() *QueryParameters {
	return &QueryParameters{}
}

type PackageInfo struct {
	Name    string
	Version string
}

// Get information about installed packages.
func Query(pkgs []string, opts ...QueryParameters) error {
	o := formOptions(opts, QueryDefault)

	args := []string{"-Q"}
	if o.Explicit {
		args = append(args, "--explicit")
	}
	if o.Deps {
		args = append(args, "--deps")
	}
	if o.Native {
		args = append(args, "--native")
	}
	if o.Foreign {
		args = append(args, "--foreign")
	}
	if o.Unrequired {
		args = append(args, "--unrequired")
	}
	if o.Groups {
		args = append(args, "--groups")
	}
	for range o.List {
		args = append(args, "-l")
	}
	for range o.Info {
		args = append(args, "-i")
	}
	for range o.Check {
		args = append(args, "-k")
	}
	if o.Upgrade {
		args = append(args, "-u")
	}
	if o.File != "" {
		args = append(args, "--file")
		args = append(args, o.File)
	}

	args = append(args, o.AdditionalParams...)
	args = append(args, pkgs...)

	cmd := exec.Command(pacman, args...)
	cmd.Stdout = o.Stdout
	cmd.Stderr = o.Stderr
	cmd.Stdin = o.Stdin

	return cmd.Run()
}

type PackageInfoFull struct {
	Name          string
	Version       string
	Description   string
	Architecture  string
	URL           string
	Licenses      string
	Groups        string
	Provides      string
	DependsOn     string
	OptionalDeps  string
	RequiredBy    string
	OptionalFor   string
	ConflictsWith string
	Replaces      string
	InstalledSize string
	Packager      string
	BuildDate     string
	InstallDate   string
	InstallReason string
	InstallScript string
	ValidatedBy   string
}

// Get info about package.
func Info(pkg string) (*PackageInfoFull, error) {
	var b bytes.Buffer
	cmd := exec.Command(pacman, "-Qi", pkg)
	cmd.Stdout = &b
	cmd.Stderr = &b

	err := cmd.Run()
	if err != nil {
		return nil, errors.New("unable to get info: " + b.String())
	}
	out := b.String()

	return &PackageInfoFull{
		Name:          parseField(out, "Name            : "),
		Version:       parseField(out, "Version         : "),
		Description:   parseField(out, "Description     : "),
		Architecture:  parseField(out, "Architecture    : "),
		URL:           parseField(out, "URL             : "),
		Licenses:      parseField(out, "Licenses        : "),
		Groups:        parseField(out, "Groups          : "),
		Provides:      parseField(out, "Provides        : "),
		DependsOn:     parseField(out, "Depends On      : "),
		OptionalDeps:  parseField(out, "Optional Deps   : "),
		RequiredBy:    parseField(out, "Required By     : "),
		OptionalFor:   parseField(out, "Optional For    : "),
		ConflictsWith: parseField(out, "Conflicts With  : "),
		Replaces:      parseField(out, "Replaces        : "),
		InstalledSize: parseField(out, "Installed Size  : "),
		Packager:      parseField(out, "Packager        : "),
		BuildDate:     parseField(out, "Build Date      : "),
		InstallDate:   parseField(out, "Install Date    : "),
		InstallReason: parseField(out, "Install Reason  : "),
		InstallScript: parseField(out, "Install Script  : "),
		ValidatedBy:   parseField(out, "Validated By    : "),
	}, nil
}

func parseField(full string, field string) string {
	splt := strings.Split(full, field)
	return strings.Split(splt[1], "\n")[0]
}

// Outdated package.
type OutdatedPackage struct {
	Name           string
	CurrentVersion string
	NewVersion     string
}

// Get information about outdated packages.
func Outdated() ([]OutdatedPackage, error) {
	var b bytes.Buffer
	cmd := exec.Command(pacman, "-Qu")
	cmd.Stdout = &b
	cmd.Stderr = &b

	err := cmd.Run()
	if err != nil {
		if b.String() == `` {
			return nil, nil
		}
		return nil, errors.New("unable to get info: " + b.String())
	}
	out := b.String()
	return parseOutdated(out), nil
}

func parseOutdated(o string) []OutdatedPackage {
	var rez []OutdatedPackage
	for _, line := range strings.Split(o, "\n") {
		if line == "" {
			break
		}
		splt := strings.Split(line, " ")
		rez = append(rez, OutdatedPackage{
			Name:           splt[0],
			CurrentVersion: splt[1],
			NewVersion:     splt[3],
		})
	}
	return rez
}

// Get raw file infor for provided package using `pacman -Qp`.
func RawFileInfo(filepath string) (string, error) {
	var b bytes.Buffer
	cmd := exec.Command("pacman", "-Qpi", filepath)
	cmd.Stdout = &b
	cmd.Stderr = &b
	err := cmd.Run()
	if err != nil {
		return ``, errors.New(b.String())
	}
	return b.String(), nil
}
