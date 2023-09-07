// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package creds

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"fmnx.su/core/pack/msgs"
)

// Get login and password from git credentials file in user directory.
func Get(protocol, addr string) (string, string, error) {
	userdir, err := os.UserHomeDir()
	if err != nil {
		return ``, ``, err
	}

	b, err := os.ReadFile(path.Join(userdir, ".git-credentials"))
	if err != nil {
		return ``, ``, err
	}

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, protocol+":") && strings.HasSuffix(line, addr) {

			splt := strings.Split(line, ":")

			if len(splt) < 2 {
				return ``, ``, errors.New("bad git credentials")
			}

			login := strings.ReplaceAll(splt[1], "//", "")
			password := strings.Split(splt[2], "@")[0]
			return login, password, nil
		}
	}

	return ``, ``, errors.New("no credentials for " + addr)
}

// Save/Update login and password in git credentials file.
func Put(protocol, addr, login, pass string) error {
	userdir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	credsfile := path.Join(userdir, ".git-credentials")
	credsline := fmt.Sprintf("%s://%s:%s@%s\n", protocol, login, pass, addr)

	b, err := os.ReadFile(credsfile)
	if errors.Is(err, os.ErrNotExist) || err == nil {
		newcontents := append(b, []byte(credsline)...)
		return os.WriteFile(credsfile, newcontents, os.ModePerm)
	}
	return err
}

// Create new pair of credentials, requires user input. Returns login and
// password as a result.
func Create(protocol, addr string, r io.Reader, w io.Writer) (string, string, error) {
	login, err := msgs.Inp("Enter login: ", w, r, false)
	if err != nil {
		return ``, ``, err
	}

	pass, err := msgs.Inp("Enter password: ", w, r, true)
	if err != nil {
		return ``, ``, err
	}

	err = Put(protocol, addr, login, pass)
	if err != nil {
		return ``, ``, err
	}

	return login, pass, nil
}
