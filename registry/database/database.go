// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package database

type Database interface {
	Save(key string, data []byte) error
	Load(key string) ([]byte, error)
	Remove(key string) error
	DbKeys() ([]string, error)
	GpgKey(owner, email string) (string, error)
}
