// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"fmnx.su/core/pack/msgs"
	"github.com/jessevdk/go-flags"
)

func getOptions[Opts any](arr []Opts) *Opts {
	if len(arr) == 1 {
		return &arr[0]
	}

	var opts Opts
	_, err := flags.NewParser(&opts, flags.IgnoreUnknown).Parse()
	if err != nil {
		fmt.Println(msgs.Err + err.Error())
		os.Exit(1)
	}

	return &opts
}

func call(cmd *exec.Cmd) error {
	var buf bytes.Buffer
	cmd.Stderr = &buf
	err := cmd.Run()
	if err != nil {
		out := strings.ReplaceAll(buf.String(), "error: ", "")
		return errors.New(strings.TrimSuffix(out, "\n"))
	}
	return nil
}
