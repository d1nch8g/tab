// Use of this code is governed by GNU General Public License.
// Official web page: https://ion.lc/core/tab
// Contact email: help@ion.lc

package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/jessevdk/go-flags"
	"ion.lc/core/tab/msgs"
	"ion.lc/core/tab/tab"
)

var opts struct {
	Help    bool `long:"help" short:"h"`
	Version bool `long:"version" short:"v"`

	Query  bool `short:"Q" long:"query"`
	Remove bool `short:"R" long:"remove"`
	Sync   bool `short:"S" long:"sync"`
	Push   bool `short:"P" long:"push"`
	Build  bool `short:"B" long:"build"`
}

var help = `Simplified version of pacman

operations:
	tab {-S --sync}   [options] [(registry)/(owner)/package(s)]
	tab {-P --push}   [options] [(registry)/(owner)/package(s)]
	tab {-R --remove} [options] [(registry)/(owner)/package(s)]
	tab {-B --build}  [options] [git/repository(s)]
	tab {-Q --query}  [options] [package(s)]

use 'pack {-h --help}' with an operation for available options`

var version = `             Tab - package manager
            Copyright  (C) 2023 ION
     
  This program may be freely redistributed under
   the terms of the GNU General Public License.
       Web page: https://ion.lc/core/tab
 
                Version: 0.2.0`

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

	switch {
	case opts.Sync && opts.Help:
		fmt.Println(tab.SyncHelp)
		return nil

	case opts.Sync:
		return tab.Sync(args())

	case opts.Push && opts.Help:
		fmt.Println(tab.PushHelp)
		return nil

	case opts.Push:
		return tab.Push(args())

	case opts.Remove && opts.Help:
		fmt.Println(tab.RemoveHelp)
		return nil

	case opts.Remove:
		return tab.Remove(args())

	case opts.Query && opts.Help:
		fmt.Println(tab.QueryHelp)
		return nil

	case opts.Query:
		return tab.Query(args())

	case opts.Build && opts.Help:
		fmt.Println(tab.BuildHelp)
		return nil

	case opts.Build:
		return tab.Build(args())

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
