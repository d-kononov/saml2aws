module github.com/d-kononov/saml2aws/v2

go 1.15

require (
	github.com/99designs/keyring v1.1.6
	github.com/AlecAivazis/survey/v2 v2.3.6
	github.com/Azure/go-ntlmssp v0.0.0-20180416175057-4b934ac9dad3
	github.com/PuerkitoBio/goquery v1.5.1
	github.com/alecthomas/kingpin v2.2.6+incompatible
	github.com/avast/retry-go v2.6.0+incompatible
	github.com/aws/aws-sdk-go v1.40.9
	github.com/beevik/etree v1.0.1
	github.com/danieljoos/wincred v1.1.0
	github.com/google/uuid v1.2.0
	github.com/keybase/go-keychain v0.0.0-20190712205309-48d3d31d256d
	github.com/marshallbrekka/go-u2fhost v0.0.0-20210111072507-3ccdec8c8105
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mxschmitt/playwright-go v0.1400.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.8.1
	github.com/versent/saml2aws/v2 v2.33.0
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
	gopkg.in/ini.v1 v1.62.0
	sigs.k8s.io/aws-iam-authenticator v0.5.3
)

// replace github.com/keybase/go-keychain => github.com/wolfeidau/go-keychain v0.0.0-20210215232950-1e19148f864f
