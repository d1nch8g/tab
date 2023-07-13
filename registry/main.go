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
)

var cfg struct {
	Name    string `env:"PACK_REGISTRY_NAME"`
	Port    string `env:"PACK_REGISTRY_PORT"`
	DbDir   string `env:"PACK_REGISTRY_DB_DIR"`
	GpgDir  string `env:"PACK_REGISTRY_GPG_DIR"`
	TlsCert string `env:"PACK_REGISTRY_TLS_CERT"`
	TlsKey  string `env:"PACK_REGISTRY_TLS_KEY"`
}

func main() {
	err := run()
	if err != nil {
		fmt.Println(msgs.Err + err.Error())
		os.Exit(1)
	}
}

func run() error {
	err := env.Parse(&cfg)
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
