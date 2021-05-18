// +build darwin,cgo

package commands

import (
	"github.com/d-kononov/saml2aws/v2/helper/credentials"
	"github.com/d-kononov/saml2aws/v2/helper/osxkeychain"
)

func init() {
	credentials.CurrentHelper = &osxkeychain.Osxkeychain{}
}
