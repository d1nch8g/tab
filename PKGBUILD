# Maintainer: Danila Fominykh <d1nch8g@ion.lc>

pkgname="tab"
pkgver="0.2.2"
pkgrel="1"
pkgdesc="Decentralized package manager based on pacman."
arch=("x86_64")
url="https://ion.lc/core/tab"
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
  install -Dm755 p $pkgdir/usr/bin/tab
}
