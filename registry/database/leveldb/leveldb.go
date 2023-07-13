// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package leveldb

import (
	"bytes"
	"os"
	"path"

	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/crypto/openpgp"
)

// Database parameters.
type DatabaseParams struct {
	// Directory, where keys related to database will be stored.
	GpgKeyDirectory string
	// Directory, where leveldb files will be stored.
	LevelDbDirectory string
}

// Function will read directory with GPG keys, save them to database and
// initilize database instance in database directory.
func Get(p *DatabaseParams) (*Database, error) {
	userspace, err := os.ReadDir(p.GpgKeyDirectory)
	if err != nil {
		return nil, err
	}

	ldb, err := leveldb.OpenFile(p.LevelDbDirectory, nil)
	if err != nil {
		return nil, err
	}

	for _, userdir := range userspace {
		if userdir.IsDir() {
			userdirpath := path.Join(p.GpgKeyDirectory, userdir.Name())
			keyDes, err := os.ReadDir(userdirpath)
			if err != nil {
				return nil, err
			}

			for _, keyDe := range keyDes {
				keypath := path.Join(userdirpath, keyDe.Name())
				keydata, err := os.ReadFile(keypath)
				if err != nil {
					return nil, err
				}

				_, err = openpgp.ReadArmoredKeyRing(bytes.NewReader(keydata))
				if err != nil {
					return nil, err
				}

				keyentry := []byte(userdir.Name() + "." + keyDe.Name())
				err = ldb.Put(keyentry, keydata, nil)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return &Database{
		DB: ldb,
	}, nil
}

type Database struct {
	*leveldb.DB
}

func (d *Database) Save(key string, data []byte) error {
	return d.DB.Put([]byte(key), data, nil)
}

func (d *Database) Load(key string) ([]byte, error) {
	return d.DB.Get([]byte(key), nil)
}

func (d *Database) Remove(key string) error {
	return d.DB.Delete([]byte(key), nil)
}

func (d *Database) DbKeys() ([]string, error) {
	var keys []string
	iter := d.DB.NewIterator(nil, nil)
	for iter.Next() {
		keys = append(keys, string(iter.Key()))
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (d *Database) GpgKey(owner, email string) (string, error) {
	keydata, err := d.DB.Get([]byte(owner+"."+email), nil)
	if err != nil {
		return ``, err
	}
	return string(keydata), nil
}
