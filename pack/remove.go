// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"fmnx.su/core/pack/msgs"
	"fmnx.su/core/pack/pacman"
)

type RemoveParameters struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	// Ask for confirmation when deleting package.
	Confirm bool
	// Leave package dependencies in the system (removed by default).
	Norecursive bool
	// Leave package configs in the system (removed by default).
	Nocfgs bool
	// Remove packages and all packages that depend on them.
	Cascade bool
	// Custom distribution name that will be used for package deletion.
	Distro string
	// Use insecure connection for remote deletions.
	Insecure bool
	// Set custom architectures for deletion.
	Arch string
}

func removeDefault() *RemoveParameters {
	return &RemoveParameters{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
		Distro: "archlinux",
		Arch:   "x86_64",
	}
}

func Remove(args []string, prms ...RemoveParameters) error {
	p := formOptions(prms, removeDefault)

	local, remote := splitRemoved(args)

	if len(local) > 0 {
		err := pacman.RemoveList(local, pacman.RemoveParameters{
			Sudo:        true,
			NoConfirm:   !p.Confirm,
			Recursive:   !p.Norecursive,
			WithConfigs: !p.Nocfgs,
			Cascade:     p.Cascade,
			Stdout:      p.Stdout,
			Stderr:      p.Stderr,
			Stdin:       p.Stdin,
		})
		if err != nil {
			return err
		}
	}

	if len(remote) > 0 {
		email, err := gnupgEmail()
		if err != nil {
			return err
		}
		if p.Distro == "" {
			p.Distro = "archlinux"
		}
		msgs.Amsg(p.Stdout, "Removing remote packages as "+email)
		for i, pkg := range remote {
			msgs.Smsg(p.Stdout, "Removing "+pkg, i+1, len(remote))
			err := rmRemote(p, pkg, email)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Splits packages that will be removed locally and on remote.
func splitRemoved(pkgs []string) ([]string, []string) {
	var local []string
	var remote []string
	for _, pkg := range pkgs {
		if strings.Contains(pkg, "/") {
			remote = append(remote, pkg)
			continue
		}
		local = append(local, pkg)
	}
	return local, remote
}

// Get remote, owner, target and version from input arguement.
func splitPkg(pkg string) (string, string, string, string, error) {
	splt := strings.Split(pkg, "/")
	if len(splt) == 2 {
		pkg, ver, err := splitVer(splt[1])
		return splt[0], ``, pkg, ver, err
	}
	pkg, ver, err := splitVer(splt[2])
	return splt[0], splt[1], pkg, ver, err
}

func splitVer(pkg string) (string, string, error) {
	splt := strings.Split(pkg, "@")
	if len(splt) != 2 {
		return "", "", fmt.Errorf("unable to eject version from %s", pkg)
	}
	return splt[0], splt[1], nil
}

// Function that will be used to remove remote package.
func rmRemote(p *RemoveParameters, pkg, email string) error {
	t := time.Now().Format(time.RFC3339)

	remote, owner, target, version, err := splitPkg(pkg)
	if err != nil {
		return err
	}

	err = os.WriteFile("packdel", []byte(t+owner+target), os.ModePerm)
	if err != nil {
		return err
	}

	var errbuf bytes.Buffer
	cmd := exec.Command("gpg", "--detach-sign", "packdel")
	cmd.Stdout = p.Stdout
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		return errors.Join(errors.New(errbuf.String()), err)
	}

	err = os.RemoveAll("packdel")
	if err != nil {
		return err
	}

	signature, err := os.Open("packdel.sig")
	if err != nil {
		return err
	}

	prfx := "https://"
	if p.Insecure {
		prfx = "http://"
	}

	req, err := http.NewRequest(
		http.MethodDelete,
		prfx+path.Join(remote, "api/packages", owner, "arch/remove"),
		signature,
	)
	if err != nil {
		return err
	}

	req.Header.Add("email", email)
	req.Header.Add("distro", p.Distro)
	req.Header.Add("target", target)
	req.Header.Add("time", t)
	req.Header.Add("version", version)
	req.Header.Add("arch", p.Arch)

	var client http.Client
	resp, err := client.Do(req)
	err = errors.Join(os.RemoveAll("packdel.sig"), err)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Join(err, errors.New(resp.Status))
		}
		return fmt.Errorf("%s %s", resp.Status, string(b))
	}
	return nil
}
