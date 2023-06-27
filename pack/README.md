<h2 align="center">Pack</h2>

This library provides API's for pack specific operations, that can't be done with pacman.
It has 4 main functions:

- `Build()` - for package builds
- `Open()` - starts pack registry registry
- `Push()` - pushes packages to pack registry
- `Sync()` - syncronizes packages with pack registries

Examples:

```go

import "fmnx.su/core/pack/pack"

func main() {
    err := pack.Sync(args(), pack.SyncParameters{
        Quick:     true,
        Refresh:   true,
        Stdout:    os.Stdout,
        Stderr:    os.Stderr,
        Stdin:     os.Stdin,
        ...
    })
}

```

```go

import "fmnx.su/core/pack/pack"

func main() {
    err := pack.Push(args(), pack.SyncParameters{
        Directory: opts.Dir,
        Protocol:  opts.Protocol,
        Endpoint:  opts.Endpoint,
        ...
    })
}

```

```go

import "fmnx.su/core/pack/pack"

func main() {
    err := pack.Open(args(), pack.SyncParameters{
        Stdout:   os.Stdout,
        Stderr:   os.Stderr,
        Stdin:    os.Stdin,
        Endpoint: "/custom/endpoint",
        Dir:      "/pkg/cache/dir",
        Name:     "domain.su",
        Port:     "80",
        Cert:     "/path/to/cert.pem",
        Key:      "/path/to/key.pem",
        GpgDir:   "/custom/gpg/dir",
    })
}

```
