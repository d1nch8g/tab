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
	Gpg    bool `short:"G" long:"gpg"`
	Tmpl   bool `short:"T" long:"tmpl"`
}

var help = `Simplified version of pacman

operations:
	pack {-S --sync}   [options] [(registry)/(owner)/package(s)]
	pack {-P --push}   [options] [(registry)/(owner)/package(s)]
	pack {-R --remove} [options] [(registry)/(owner)/package(s)]
	pack {-B --build}  [options] [git/repository(s)]
	pack {-Q --query}  [options] [package(s)]
	pack {-G --gpg}    [options] [args]
	pack {-T --tmpl}   [options] [args]

use 'pack {-h --help}' with an operation for available options`

var version = `             Pack - package manager.
          Copyright (C) 2023 FMNX team
     
  This program may be freely redistributed under
   the terms of the GNU General Public License.
       Web page: https://fmnx.su/core/pack
 
                 Version: 0.1.7`

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
	_, err := flags.NewParser(&opts, flags.IgnoreUnknown).Parse()
	if err != nil {
		return err
	}
	RemoveCapitalArgs()

	switch {
	case opts.Sync && opts.Help:
		fmt.Println(pack.SyncHelp)
		return nil

	case opts.Sync:
		return pack.Sync(args())

	case opts.Push && opts.Help:
		fmt.Println(pack.PushHelp)
		return nil

	case opts.Push:
		return pack.Push(args())

	case opts.Remove && opts.Help:
		fmt.Println(pack.RemoveHelp)
		return nil

	case opts.Remove:
		return pack.Remove(args())

	case opts.Query && opts.Help:
		fmt.Println(pack.QueryHelp)
		return nil

	case opts.Query:
		return pack.Query(args())

	case opts.Build && opts.Help:
		fmt.Println(pack.BuildHelp)
		return nil

	case opts.Build:
		return pack.Build(args())

	case opts.Gpg && opts.Help:
		fmt.Println(pack.GpgHelp)
		return nil

	case opts.Gpg:
		return pack.Gpg(args())

	case opts.Tmpl && opts.Help:
		fmt.Println(pack.TmplHelp)
		return nil

	case opts.Tmpl:
		return pack.Tmpl(args())

	case opts.Version:
		fmt.Println(version)
		return nil

	case opts.Help:
		fmt.Println(help)
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

// TODO: remove lated when bug with unknown args fixed.
func RemoveCapitalArgs() {
	var newargs []string
	for _, v := range os.Args {
		if strings.HasPrefix(v, "-") {
			rootargs := []string{"Q", "R", "S", "P", "B", "G", "T"}
			for _, letter := range rootargs {
				v = strings.Replace(v, letter, "", 1)
			}
		}
		newargs = append(newargs, v)
	}
	os.Args = newargs
}
