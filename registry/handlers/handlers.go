// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package handlers

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"fmnx.su/core/pack/registry/context"
	"fmnx.su/core/pack/registry/metadata"
	"fmnx.su/core/pack/registry/service"
)

// Function used to add new packages to registry.
func Push(ctx *context.Context) {
	var (
		owner    = ctx.Params("username")
		filename = ctx.Req.Header.Get("filename")
		email    = ctx.Req.Header.Get("email")
		distro   = ctx.Req.Header.Get("distro")
		sendtime = ctx.Req.Header.Get("time")
		pkgsign  = ctx.Req.Header.Get("pkgsign")
		metasign = ctx.Req.Header.Get("metasign")
	)

	// Decoding package signature.
	sigdata, err := hex.DecodeString(pkgsign)
	if err != nil {
		apiError(ctx, http.StatusBadRequest, err)
		return
	}

	// Read package to memory for signature validation.
	pkgdata, err := io.ReadAll(ctx.Req.Body)
	if err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}
	defer ctx.Req.Body.Close()

	// Decoding time when message was created.
	t, err := time.Parse(time.RFC3339, sendtime)
	if err != nil {
		apiError(ctx, http.StatusBadRequest, err)
		return
	}

	// Check if message is outdated.
	if time.Since(t) > time.Hour {
		apiError(ctx, http.StatusUnauthorized, "outdated message")
		return
	}

	// Decoding signature related to metadata.
	msigdata, err := hex.DecodeString(metasign)
	if err != nil {
		apiError(ctx, http.StatusBadRequest, err)
		return
	}

	// Parse metadata contained in arch package archive.
	md, err := metadata.EjectMetadata(filename, distro, ctx.Domain, pkgdata)
	if err != nil {
		apiError(ctx, http.StatusBadRequest, err)
		return
	}

	// Validating metadata signature, to ensure that operation push operation
	// is initiated by original package owner.
	err = service.ValidateSignature(ctx, &service.ValidateParams{
		ValidateData:      []byte(owner + md.Name + sendtime),
		ValidateSignature: msigdata,
		Email:             email,
		Owner:             owner,
	})
	if err != nil {
		apiError(ctx, http.StatusUnauthorized, err)
		return
	}

	// Validate package signature with any of user's GnuPG keys.
	err = service.ValidateSignature(ctx, &service.ValidateParams{
		ValidateData:      pkgdata,
		ValidateSignature: sigdata,
		Email:             email,
		Owner:             owner,
	})
	if err != nil {
		apiError(ctx, http.StatusUnauthorized, err)
		return
	}

	// Save file related to arch package.
	err = service.SaveFile(ctx, &service.FileParams{
		Owner:    owner,
		Distro:   distro,
		Data:     pkgdata,
		Filename: filename,
	})
	if err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Save file related to arch package signature.
	err = service.SaveFile(ctx, &service.FileParams{
		Owner:    owner,
		Distro:   distro,
		Data:     sigdata,
		Filename: filename + ".sig",
	})
	if err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Save file related to arch package database description.
	err = service.SaveFile(ctx, &service.FileParams{
		Owner:    owner,
		Distro:   distro,
		Data:     []byte(md.GetDbDesc()),
		Filename: filename + ".desc",
	})
	if err != nil {
		apiError(ctx, http.StatusUnauthorized, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// Get file from arch package registry.
func Get(ctx *context.Context) {
	var (
		file   = ctx.Params("file")
		owner  = ctx.Params("username")
		distro = ctx.Params("distro")
		arch   = ctx.Params("arch")
	)

	// Packages are stored in different way from pacman databases, and loaded
	// with LoadPackageFile function.
	if strings.HasSuffix(file, "tar.zst") || strings.HasSuffix(file, "zst.sig") {
		pkgdata, err := service.LoadFile(ctx, &service.FileParams{
			Owner:    owner,
			Distro:   distro,
			Filename: file,
		})
		if err != nil {
			apiError(ctx, http.StatusNotFound, err)
			return
		}

		name := metadata.Join(distro, owner, file)
		http.ServeContent(ctx.Resp, ctx.Req, name, time.Now(), bytes.NewReader(pkgdata))
		return
	}

	// Pacman databases is not stored in gitea's storage, it is created for
	// incoming request and cached.
	if strings.HasSuffix(file, ".db.tar.gz") || strings.HasSuffix(file, ".db") {
		db, err := service.CreatePacmanDb(ctx, owner, arch, distro)
		if err != nil {
			apiError(ctx, http.StatusInternalServerError, err)
			return
		}

		name := metadata.Join(distro, owner, file)
		http.ServeContent(ctx.Resp, ctx.Req, name, time.Now(), bytes.NewReader(db))
		return
	}

	ctx.Resp.WriteHeader(http.StatusNotFound)
}

// Remove specific package version, related files and pacman database entry.
func Remove(ctx *context.Context) {
	var (
		owner   = ctx.Params("username")
		email   = ctx.Req.Header.Get("email")
		target  = ctx.Req.Header.Get("target")
		stime   = ctx.Req.Header.Get("time")
		distro  = ctx.Req.Header.Get("distro")
		arch    = ctx.Req.Header.Get("arch")
		version = ctx.Req.Header.Get("version")
	)

	// Parse sent time and check if it is within last minute.
	t, err := time.Parse(time.RFC3339, stime)
	if err != nil {
		apiError(ctx, http.StatusBadRequest, err)
		return
	}

	if time.Since(t) > time.Minute {
		apiError(ctx, http.StatusUnauthorized, "outdated message")
		return
	}

	// Read signature data from request body.
	sigdata, err := io.ReadAll(ctx.Req.Body)
	if err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}
	defer ctx.Req.Body.Close()

	// Validate package signature with any of user's GnuPG keys.
	err = service.ValidateSignature(ctx, &service.ValidateParams{
		ValidateData:      []byte(owner + target + stime),
		ValidateSignature: sigdata,
		Email:             email,
		Owner:             owner,
	})
	if err != nil {
		apiError(ctx, http.StatusUnauthorized, err)
		return
	}

	// Remove package files and pacman database entry.
	err = service.RemoveFile(ctx, &service.FileParams{
		Owner:    owner,
		Distro:   distro,
		Filename: fmt.Sprintf("%s-%s-%s.pkg.tar.zst", target, version, arch),
	})
	if err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Resp.WriteHeader(http.StatusOK)
}

func apiError(ctx *context.Context, status int, obj interface{}) {
	log.Printf("api error: %v", obj)
	ctx.Resp.Write([]byte(fmt.Sprintf("%v", obj)))
}
