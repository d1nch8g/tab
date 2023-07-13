// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package service

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"fmnx.su/core/pack/registry/context"
	"fmnx.su/core/pack/registry/metadata"
	"github.com/keybase/go-crypto/openpgp"
)

type ValidateParams struct {
	ValidateData      []byte
	ValidateSignature []byte
	Email             string
	Owner             string
}

func ValidateSignature(ctx *context.Context, p *ValidateParams) error {
	key, err := ctx.Db.GpgKey(p.Owner, p.Email)
	if err != nil {
		return err
	}

	kr, err := openpgp.ReadArmoredKeyRing(strings.NewReader(key))
	if err != nil {
		return fmt.Errorf("unable to get keys for %s: %v", p.Owner, err)
	}
	_, err = openpgp.CheckDetachedSignature(
		kr, bytes.NewReader(p.ValidateData),
		bytes.NewReader(p.ValidateSignature),
	)
	return err
}

type FileParams struct {
	Data     []byte
	Owner    string
	Filename string
	Distro   string
}

func SaveFile(ctx *context.Context, p *FileParams) error {
	return ctx.Db.Save(metadata.Join(p.Distro, p.Owner, p.Filename), p.Data)
}

func LoadFile(ctx *context.Context, p *FileParams) ([]byte, error) {
	return ctx.Db.Load(metadata.Join(p.Distro, p.Owner, p.Filename))
}

func RemoveFile(ctx *context.Context, p *FileParams) error {
	filename := metadata.Join(p.Distro, p.Owner, p.Filename)
	return errors.Join(
		ctx.Db.Remove(filename),
		ctx.Db.Remove(filename+".sig"),
		ctx.Db.Remove(filename+".desc"),
	)
}

func CreatePacmanDb(ctx *context.Context, owner, arch, distro string) ([]byte, error) {
	keys, err := ctx.Db.DbKeys()
	if err != nil {
		return nil, err
	}

	var descs []*metadata.DescParams
	for _, key := range keys {
		if strings.HasPrefix(key, distro+"."+arch) && strings.HasSuffix(key, ".desc") {
			name, ver, err := metadata.EjectNameParameters(key)
			if err != nil {
				return nil, err
			}

			data, err := ctx.Db.Load(key)
			if err != nil {
				return nil, err
			}

			descs = append(descs, &metadata.DescParams{
				Name:    name,
				Version: ver,
				Desc:    string(data),
			})
		}
	}

	return metadata.CreatePacmanDb(descs)
}
