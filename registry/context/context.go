// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package context

import (
	"context"
	"net/http"

	"fmnx.su/core/pack/registry/database"
	"github.com/gorilla/mux"
)

type Context struct {
	Ctx    context.Context
	Req    *http.Request
	Resp   http.ResponseWriter
	Db     database.Database
	Domain string
}

func (c *Context) Params(key string) string {
	return mux.Vars(c.Req)[key]
}

func (c *Context) Status(status int) {
	c.Resp.WriteHeader(status)
}
