# Maintainer: Danila Fominykh <dancheg97@fmnx.su>

pkgname="pack"
pkgver="0.1.9"
pkgrel="1"
pkgdesc="Decentralized package manager based on pacman."
arch=("x86_64")
url="https://fmnx.su/core/pack"
license=("GPL")
depends=("pacman")
makedepends=("go")
optdepends=(
  "git: build remote repositories"
  "sudo: privilege elevation"
  "doas: privilege elevation"
)

build() {
  go build -o p ../.
}

package() {
  install -Dm755 p $pkgdir/usr/bin/pack
}
