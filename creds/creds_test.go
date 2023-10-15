// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://ion.lc/core/tab
// Contact email: help@ion.lc

package creds

import (
	"os"
	"path"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestGetNormal(t *testing.T) {
	creds := substitute("https://user:password@example.com\n")
	name, pass, err := Get("https", "example.com")
	save(creds)

	assert.NoError(t, err)
	assert.Equal(t, "user", name)
	assert.Equal(t, "password", pass)
}

func TestGetNotExist(t *testing.T) {
	creds := substitute("")
	_, _, err := Get("https", "example.com")
	save(creds)

	assert.Equal(t, ErrNoCreds, err)
}

func TestPut(t *testing.T) {
	creds := substitute("")
	err1 := Put("https", "example.com", "user", "password")
	name, pass, err2 := Get("https", "example.com")
	save(creds)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, "user", name)
	assert.Equal(t, "password", pass)
}

func TestCreate(t *testing.T) {
	// TODO: add test for create function

	// var in bytes.Buffer
	// var out bytes.Buffer

	// in.Write([]byte("name"))
	// in.Write([]byte{'\n'})

	// in.Write([]byte("password"))
	// in.Write([]byte{'\n'})

	// creds := substitute("")
	// name, pass, err := Create("https", "example.com", &in, &out)
	// save(creds)

	// assert.NoError(t, err)
	// assert.Equal(t, "name", name)
	// assert.Equal(t, "password", pass)
}

// Helper funciton to temporarily substitute git credentials in home dir.
func substitute(fake string) string {
	save(fake)

	userdir, err := os.UserHomeDir()
	check(err)

	b, err := os.ReadFile(path.Join(userdir, ".git-credentials"))
	check(err)

	return string(b)
}

// Helper funciton to temporarily substitute git credentials in home dir.
func save(creds string) {
	userdir, err := os.UserHomeDir()
	check(err)

	dir := path.Join(userdir, ".git-credentials")
	err = os.WriteFile(dir, []byte(creds), os.ModePerm)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
