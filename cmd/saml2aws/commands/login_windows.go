package commands

import (
	"github.com/d-kononov/saml2aws/v2/helper/credentials"
	"github.com/d-kononov/saml2aws/v2/helper/wincred"
)

func init() {
	credentials.CurrentHelper = &wincred.Wincred{}
}
