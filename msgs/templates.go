// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package msgs

const PKGBUILDflutter = `# Maintainer: %s

pkgname="%s"
pkgdesc="Useful description"
pkgver="1"
pkgrel="1"
arch=('x86_64')
url="https://example.com/owner/repo"
depends=()
makedepends=(
  "flutter"
  "clang"
  "cmake"
)

build() {
  cd ..
  flutter build linux
}

package() {
  cd ..
  install -Dm755 %s.sh $pkgdir/usr/bin/%s
  install -Dm755 %s.desktop $pkgdir/usr/share/applications/%s.desktop
  install -Dm755 assets/%s.png $pkgdir/usr/share/icons/hicolor/512x512/apps/%s.png
  cd build/linux/x64/release/bundle
  find . -type f -exec install -Dm755 {} $pkgdir/usr/share/%s/{} \;
}
`

const PKGBUILDgocli = `# Maintainer: %s

pkgname="%s"
pkgdesc="Useful description"
pkgver="1"
pkgrel="1"
arch=('x86_64')
url="https://example.com/owner/repo"
depends=()
makedepends=(
  "go"
)

build() {
  cd ..
  go build -o src/p .
}

package() {
  install -Dm755 p $pkgdir/usr/bin/%s
}
`

const Desktop = `[Desktop Entry]
Name=Awesome Application
GenericName=Awesome Application
Comment=Awesome Application
Exec=/usr/share/%s/%s
WMClass=%s
Icon=/usr/share/%s/data/flutter_assets/assets/%s.png
Type=Application
`

const ShFile = `#!/usr/bin/env sh
exec /usr/share/%s/%s`
