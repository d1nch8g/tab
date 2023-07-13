// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"fmnx.su/core/pack/msgs"
	"fmnx.su/core/pack/pack"
	"fmnx.su/core/pack/pacman"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Help    bool `long:"help" short:"h"`
	Version bool `long:"version" short:"v"`

	// Root options.
	Query  bool `short:"Q" long:"query"`
	Remove bool `short:"R" long:"remove"`
	Sync   bool `short:"S" long:"sync"`
	Push   bool `short:"P" long:"push"`
	Build  bool `short:"B" long:"build"`
	Assist bool `short:"A" long:"assist"`

	// Sync options.
	Quick   bool   `short:"q" long:"quick"`
	Refresh []bool `short:"y" long:"refresh"`
	Upgrade []bool `short:"u" long:"upgrade"`
	Force   bool   `short:"f" long:"force"`

	// Push options.
	Dir      string `short:"d" long:"dir" default:"/var/cache/pacman/pkg"`
	Insecure bool   `short:"w" long:"insecure"`
	Distro   string `long:"distro" default:"archlinux"`
	Endpoint string `long:"endpoint" default:"/api/packages/arch"`

	// Remove options.
	Confirm     bool   `short:"c" long:"confirm"`
	Norecursive bool   `short:"a" long:"norecursive"`
	Nocfgs      bool   `short:"j" long:"nocfgs"`
	Cascade     bool   `short:"k" long:"cascade"`
	Arch        string `long:"arch" default:"x86_64"`

	// Query options.
	Info     []bool `short:"i" long:"info"`
	List     []bool `short:"l" long:"list"`
	Outdated bool   `short:"o" long:"outdated"`

	// Build options.
	Syncbuild bool `short:"s" long:"syncbuild"`
	Rmdeps    bool `short:"r" long:"rmdeps"`
	Garbage   bool `short:"g" long:"garbage"`

	// Assist options.
	Export  bool `short:"e" long:"export"`
	Gen     bool `short:"n" long:"gen"`
	Recv    bool `long:"recv"`
	Setpkgr bool `long:"setpkgr"`
	Flutter bool `long:"flutter"`
	Gocli   bool `long:"gocli"`
}

func main() {
	err := run()
	if err != nil {
		if !strings.Contains(err.Error(), "exit status 1") {
			fmt.Println(msgs.Err + err.Error())
		}
		os.Exit(1)
	}
}

func run() error {
	_, err := flags.NewParser(&opts, flags.None).Parse()
	if err != nil {
		return err
	}

	switch {
	case opts.Sync && opts.Help:
		fmt.Println(msgs.SyncHelp)
		return nil

	case opts.Sync:
		return pack.Sync(args(), pack.SyncParameters{
			Quick:    opts.Quick,
			Refresh:  opts.Refresh,
			Upgrade:  opts.Upgrade,
			Force:    opts.Force,
			Insecure: opts.Insecure,
			Stdout:   os.Stdout,
			Stderr:   os.Stderr,
			Stdin:    os.Stdin,
		})

	case opts.Push && opts.Help:
		fmt.Println(msgs.PushHelp)
		return nil

	// TODO: when multiple architectures found in cache push them all.
	case opts.Push:
		return pack.Push(args(), pack.PushParameters{
			Stdout:    os.Stdout,
			Stderr:    os.Stderr,
			Stdin:     os.Stdin,
			Directory: opts.Dir,
			Insecure:  opts.Insecure,
			Distro:    opts.Distro,
		})

	case opts.Remove && opts.Help:
		fmt.Println(msgs.RemoveHelp)
		return nil

	case opts.Remove:
		return pack.Remove(args(), pack.RemoveParameters{
			Stdout:      os.Stdout,
			Stderr:      os.Stderr,
			Stdin:       os.Stdin,
			Confirm:     opts.Confirm,
			Norecursive: opts.Norecursive,
			Nocfgs:      opts.Nocfgs,
			Cascade:     opts.Cascade,
			Distro:      opts.Distro,
			Insecure:    opts.Insecure,
		})

	case opts.Query && opts.Help:
		fmt.Println(msgs.QueryHelp)
		return nil

	case opts.Query:
		if opts.Outdated {
			return pacman.Query(nil, pacman.QueryParameters{
				Stdout:  os.Stdout,
				Stderr:  os.Stderr,
				Stdin:   os.Stdin,
				Upgrade: true,
			})
		}
		return pacman.Query(args(), pacman.QueryParameters{
			Info:   opts.Info,
			List:   opts.List,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Stdin:  os.Stdin,
		})

	case opts.Build && opts.Help:
		fmt.Println(msgs.BuildHelp)
		return nil

	case opts.Build:
		return pack.Build(args(), pack.BuildParameters{
			Dir:       opts.Dir,
			Quick:     opts.Quick,
			Syncbuild: opts.Syncbuild,
			Rmdeps:    opts.Rmdeps,
			Garbage:   opts.Garbage,
			Stdout:    os.Stdout,
			Stderr:    os.Stderr,
			Stdin:     os.Stdin,
		})

	case opts.Assist && opts.Help:
		fmt.Println(msgs.AssistHelp)
		return nil

	case opts.Assist:
		return pack.Assist(args(), pack.AssistParameters{
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
			Stdin:   os.Stdin,
			Export:  opts.Export,
			Gen:     opts.Gen,
			Recv:    opts.Recv,
			Info:    len(opts.Info) > 0,
			Setpkgr: opts.Setpkgr,
			Flutter: opts.Flutter,
			Gocli:   opts.Gocli,
		})

	case opts.Version:
		fmt.Println(msgs.Version)
		return nil

	case opts.Help:
		fmt.Println(msgs.Help)
		return nil

	default:
		return fmt.Errorf("specify at least one root flag (pack -h)")
	}
}

// Function to get list of command line arguements. It automatically filters
// all string CLI parameters with reflect.
func args() []string {
	var arglist []string

	v := reflect.ValueOf(opts)

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Type().String() == "string" {
			short := reflect.TypeOf(opts).Field(i).Tag.Get("short")
			if short != "" {
				arglist = append(arglist, "-"+short)
			}
			long := reflect.TypeOf(opts).Field(i).Tag.Get("long")
			arglist = append(arglist, "--"+long)
		}
	}

	var filtered []string
	for i, v := range os.Args {
		if i == 0 || i == 1 {
			continue
		}
		if strings.HasPrefix(v, "-") {
			continue
		}
		var skipStringArg bool
		for _, args := range arglist {
			if os.Args[i-1] == args {
				skipStringArg = true
			}
		}
		if skipStringArg {
			continue
		}
		filtered = append(filtered, v)
	}
	return filtered
}
