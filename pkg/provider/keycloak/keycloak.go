package keycloak

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/d-kononov/saml2aws/v2/pkg/cfg"
	"github.com/d-kononov/saml2aws/v2/pkg/creds"
	"github.com/d-kononov/saml2aws/v2/pkg/prompter"
	"github.com/d-kononov/saml2aws/v2/pkg/provider"
	"github.com/pkg/errors"
)

// Client wrapper around KeyCloak.
type Client struct {
	client       *provider.HTTPClient
	silentOutput bool
	provider.ValidateBase
}

// New create a new KeyCloakClient
func New(idpAccount *cfg.IDPAccount) (*Client, error) {

	tr := provider.NewDefaultTransport(idpAccount.SkipVerify)

	client, err := provider.NewHTTPClient(tr, provider.BuildHttpClientOpts(idpAccount))
	if err != nil {
		return nil, errors.Wrap(err, "error building http client")
	}

	return &Client{
		client:       client,
		silentOutput: idpAccount.SilentOutput,
	}, nil
}

// Authenticate logs into KeyCloak and returns a SAML response
func (kc *Client) Authenticate(loginDetails *creds.LoginDetails) (string, error) {

	authSubmitURL, authForm, err := kc.getLoginForm(loginDetails)
	if err != nil {
		return "", errors.Wrap(err, "error retrieving login form from idp")
	}

	data, err := kc.postLoginForm(authSubmitURL, authForm)
	if err != nil {
		return "", fmt.Errorf("error submitting login form")
	}
	if authSubmitURL == "" {
		return "", fmt.Errorf("error submitting login form")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return "", errors.Wrap(err, "error parsing document")
	}
	if credErr := doc.Find("#input-error"); credErr.Length() > 0 {
		if err := credErr.First().Text(); err != "" {
			return "", errors.Wrap(errors.New(strings.TrimSpace(err)), "Auth failed")
		}
	}
	if containsTotpForm(doc) {
		totpSubmitURL, err := extractSubmitURL(doc)
		if err != nil {
			return "", errors.Wrap(err, "unable to locate IDP totp form submit URL")
		}
		totpCredentialId := extractCredentialId(doc)

		if loginDetails.MFAToken == "" {
			println("Please input TOTP code and hit enter")
		}
		doc, err = kc.postTotpForm(totpSubmitURL, totpCredentialId, loginDetails.MFAToken, doc)
		if otpErr := doc.Find("#input-error-otp-code"); otpErr.Length() > 0 {
			if err := otpErr.First().Text(); err != "" {
				return "", errors.Wrap(errors.New(strings.TrimSpace(err)), "Auth failed")
			}
		}
		if err != nil {
			return "", errors.Wrap(err, "error posting totp form")
		}
	}

	return extractSamlResponse(doc)
}

func (kc *Client) getLoginForm(loginDetails *creds.LoginDetails) (string, url.Values, error) {

	res, err := kc.client.Get(loginDetails.URL)
	if err != nil {
		return "", nil, errors.Wrap(err, "error retrieving form")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to build document from response")
	}

	if res.StatusCode == http.StatusUnauthorized {
		authSubmitURL, err := extractSubmitURL(doc)
		if err != nil {
			return "", nil, errors.Wrap(err, "unable to locate IDP authentication form submit URL")
		}
		loginDetails.URL = authSubmitURL
		return kc.getLoginForm(loginDetails)
	}

	authForm := url.Values{}

	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		updateKeyCloakFormData(authForm, s, loginDetails)
	})

	authSubmitURL, err := extractSubmitURL(doc)
	if err != nil {
		return "", nil, errors.Wrap(err, "unable to locate IDP authentication form submit URL")
	}

	return authSubmitURL, authForm, nil
}

func (kc *Client) postLoginForm(authSubmitURL string, authForm url.Values) ([]byte, error) {

	req, err := http.NewRequest("POST", authSubmitURL, strings.NewReader(authForm.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "error building authentication request")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := kc.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving login form")
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving body")
	}

	return data, nil
}

func (kc *Client) postTotpForm(totpSubmitURL string, totpCredentialId *string, mfaToken string, doc *goquery.Document) (*goquery.Document, error) {

	otpForm := url.Values{}

	if mfaToken == "" {
		mfaToken = prompter.RequestSecurityCode("000000", kc.silentOutput)
	}

	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		updateOTPFormData(otpForm, s, mfaToken)
	})
	if totpCredentialId != nil {
		otpForm.Add("selectedCredentialId", *totpCredentialId)
	}

	req, err := http.NewRequest("POST", totpSubmitURL, strings.NewReader(otpForm.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "error building MFA request")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := kc.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving content")
	}

	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading totp form response")
	}

	return doc, nil
}

func extractSubmitURL(doc *goquery.Document) (string, error) {

	var submitURL string

	doc.Find("form").Each(func(i int, s *goquery.Selection) {
		action, ok := s.Attr("action")
		if !ok {
			return
		}
		submitURL = action
	})

	if submitURL == "" {
		return "", fmt.Errorf("unable to locate form submit URL")
	}

	return submitURL, nil
}

func extractCredentialId(doc *goquery.Document) (credId *string) {
	doc.Find(".form-group").Each(func(i int, s *goquery.Selection) {
		if s.Find(".otp-tile").Length() > 0 {
			credIdValue, ok := s.Find("input[type='hidden']").Last().Attr("value")
			if !ok {
				return
			}
			credId = &credIdValue
			return
		}
	})
	return credId
}

func extractSamlResponse(doc *goquery.Document) (string, error) {
	var samlAssertion string
	var err error

	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		name, ok := s.Attr("name")
		if ok && name == "SAMLResponse" {
			val, ok := s.Attr("value")
			if !ok {
				err = errors.New("unable to locate saml assertion value")
			}
			samlAssertion = val
		}
	})
	if err != nil {
		return "", errors.New("unable to locate saml assertion value")
	}

	if samlAssertion == "" {
		return "", errors.New("unable to locate saml response field")
	}

	return samlAssertion, nil
}

func containsTotpForm(doc *goquery.Document) bool {
	// search totp field at Keycloak < 8.0.1
	totpIndex := doc.Find("input#totp").Index()

	if totpIndex != -1 {
		return true
	}

	// search otp field at Keycloak >= 8.0.1
	totpIndex = doc.Find("input#otp").Index()

	return totpIndex != -1
}

func updateKeyCloakFormData(authForm url.Values, s *goquery.Selection, user *creds.LoginDetails) {
	name, ok := s.Attr("name")
	// log.Printf("name = %s ok = %v", name, ok)
	if !ok {
		return
	}
	lname := strings.ToLower(name)
	if strings.Contains(lname, "username") {
		authForm.Add(name, user.Username)
	} else if strings.Contains(lname, "password") {
		authForm.Add(name, user.Password)
	} else {
		// pass through any hidden fields
		val, ok := s.Attr("value")
		if !ok {
			return
		}
		authForm.Add(name, val)
	}
}

func updateOTPFormData(otpForm url.Values, s *goquery.Selection, token string) {
	name, ok := s.Attr("name")
	// log.Printf("name = %s ok = %v", name, ok)
	if !ok {
		return
	}

	lname := strings.ToLower(name)
	// search otp field at Keycloak >= 8.0.1
	if strings.Contains(lname, "totp") {
		otpForm.Add(name, token)
	} else if strings.Contains(lname, "otp") {
		otpForm.Add(name, token)
	}

}
