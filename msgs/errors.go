// Use of this code is governed by GNU General Public License.
// Official web page: https://ion.lc/core/tab
// Contact email: help@ion.lc

package msgs

import "github.com/fatih/color"

var Err = color.New(color.Bold, color.FgHiRed).Sprintf("error: ")

var ErrGnuPGprivkeyNotFound = `GnuPG private key not found.
It is required for package signing, run:

1) Generate a key:
` + color.YellowString("gpg --gen-key") + `

2) Get KEY-ID, paste it to next command:
` + color.YellowString("gpg -k") + `

3) Send it to key GPG server:
` + color.YellowString("gpg --send-keys KEY-ID") + `

4) Edit ` + color.BlueString("PACKAGER") + ` variable in ` + color.CyanString("/etc/makepkg.conf") + `
Name and email should match with name and email in GnuPG authority for pack to work properly.
`

var ErrNoPackager = `packager not found.

Add ` + color.BlueString("PACKAGER") + ` variable matching your GnuPG authority in ` + color.CyanString("/etc/makepkg.conf") + `

` + color.BlueString("PACKAGER") + `=` + color.HiGreenString("%s")

var ErrPackagerMissmatch = `signer and packager identities are different.

Your GnuPG authority: ` + color.HiGreenString("%s") + `
Your makepkg authority: ` + color.HiCyanString("%s") + `

Authority you defined in GnuPG is not matching with ` + color.BlueString("PACKAGER") + ` variable in ` + color.CyanString("/etc/makepkg.conf")

var ErrGitUserMissmatch = `git and gnupg identities are different.

Your GnuPG authority: ` + color.HiGreenString("%s") + `
Your Git authority: ` + color.HiCyanString("%s") + `

Make sure, that ` + color.HiRedString("user.name") + ` and ` + color.HiRedString("user.email") + ` in ` + color.CyanString("~/.gitconfig") + ` are matching GnuPG name and email.`

var ErrGitSignKeyMissmatch = `git signing key does not exist in GnuPG.

Git signing key: ` + color.HiCyanString("%s") + `

You can check your gpg identity with command: ` + color.BlueString("gpg -K") + `

Make sure, that ` + color.HiRedString("user.signingkey") + ` matches key in ` + color.CyanString("~/.gitconfig")
