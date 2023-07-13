// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package main

import (
	"fmt"
	"net/http"
	"os"

	"fmnx.su/core/pack/msgs"
	"fmnx.su/core/pack/registry/context"
	"fmnx.su/core/pack/registry/database/leveldb"
	"fmnx.su/core/pack/registry/handlers"
	"github.com/caarlos0/env/v9"
	"github.com/jessevdk/go-flags"
)

var cfg struct {
	Help    bool `long:"help" short:"h"`
	Version bool `long:"version" short:"v"`

	Name    string `env:"PACK_REGISTRY_NAME" short:"n" long:"name"`
	Port    string `env:"PACK_REGISTRY_PORT" short:"p" long:"port"`
	DbDir   string `env:"PACK_REGISTRY_DB_DIR" short:"d" long:"dbdir"`
	GpgDir  string `env:"PACK_REGISTRY_GPG_DIR" short:"g" long:"gpgdir"`
	TlsCert string `env:"PACK_REGISTRY_TLS_CERT" short:"c" long:"cert"`
	TlsKey  string `env:"PACK_REGISTRY_TLS_KEY" short:"k" long:"key"`
}

func main() {
	err := run()
	if err != nil {
		fmt.Println(msgs.Err + err.Error())
		os.Exit(1)
	}
}

func run() error {
	_, err := flags.NewParser(&cfg, flags.None).Parse()
	if err != nil {
		return err
	}

	err = env.Parse(&cfg)
	if err != nil {
		return err
	}

	db, err := leveldb.Get(&leveldb.DatabaseParams{
		GpgKeyDirectory:  cfg.GpgDir,
		LevelDbDirectory: cfg.DbDir,
	})
	if err != nil {
		return err
	}

	m := http.NewServeMux()

	m.HandleFunc(
		"/api/packages/{user}/arch/push",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.Push(&context.Context{
				Ctx:    r.Context(),
				Req:    r,
				Resp:   w,
				Db:     db,
				Domain: cfg.Name,
			})
		},
	)

	m.HandleFunc(
		"/api/packages/{user}/arch/delete",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.Push(&context.Context{
				Ctx:    r.Context(),
				Req:    r,
				Resp:   w,
				Db:     db,
				Domain: cfg.Name,
			})
		},
	)

	m.HandleFunc(
		"/api/packages/{user}/arch/{distro}/{arch}/{file}",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.Push(&context.Context{
				Ctx:    r.Context(),
				Req:    r,
				Resp:   w,
				Db:     db,
				Domain: cfg.Name,
			})
		},
	)

	server := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: m,
	}

	if cfg.TlsKey != "" && cfg.TlsCert != "" {
		return server.ListenAndServeTLS(cfg.TlsCert, cfg.TlsKey)
	}
	return server.ListenAndServe()
}

const Help = `Simplified version of pacman

operations:
	pack {-S --sync}   [options] [(registry)/(owner)/package(s)]
	pack {-P --push}   [options] [(registry)/(owner)/package(s)]
	pack {-R --remove} [options] [(registry)/(owner)/package(s)]
	pack {-Q --query}  [options] [(registry)/(owner)/package(s)]
	pack {-B --build}  [options] [git/repository(s)]
	pack {-A --assist} [options] [args]

use 'pack {-h --help}' with an operation for available options`

const Version = `         Pack registry for arch packages
          Copyright (C) 2023 FMNX team
     
  This program may be freely redistributed under
   the terms of the GNU General Public License.
       Web page: https://fmnx.su/core/pack
 
                 Version: 0.1.2`
