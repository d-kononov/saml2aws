package commands

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/d-kononov/saml2aws/v2/pkg/awsconfig"
	"github.com/d-kononov/saml2aws/v2/pkg/flags"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"log"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
	"time"
)

type tokenOutput struct {
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	Status     struct {
		ExpirationTimestamp string `json:"expirationTimestamp"`
		Token               string `json:"token"`
	} `json:"status"`
}

func EksToken(configFlags *flags.LoginExecFlags) error {
	logger := logrus.WithField("command", "eks-token")

	account, err := buildIdpAccount(configFlags)
	if err != nil {
		return errors.Wrap(err, "error building login details")
	}

	sharedCreds := awsconfig.NewSharedCredentials(account.Profile, account.CredentialsFile)

	logger.Debug("check if Creds Exist")

	// this checks if the credentials file has been created yet
	exist, err := sharedCreds.CredsExists()
	if err != nil {
		return errors.Wrap(err, "error loading credentials")
	}
	if !exist {
		log.Println("unable to load credentials, login required to create them")
		return nil
	}

	awsCreds, err := sharedCreds.Load()
	if err != nil {
		return errors.Wrap(err, "error loading credentials")
	}

	if awsCreds.Expires.Sub(time.Now()) < 0 {
		if err := Login(configFlags); err != nil {
			return errors.Wrap(err, "error logging in")
		}
		awsCreds, err = sharedCreds.Load()
		if err != nil {
			return errors.Wrap(err, "error loading credentials")
		}
	}

	gen, err := token.NewGenerator(true, false)
	if err != nil {
		return errors.Wrap(err, "failed to create token generator")
	}

	sess, err := session.NewSession(&aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(true),
		Region:                        aws.String(awsCreds.Region),
		Credentials:                   credentials.NewStaticCredentials(awsCreds.AWSAccessKey, awsCreds.AWSSecretKey, awsCreds.AWSSessionToken),
	})
	if err != nil {
		return errors.Wrap(err, "failed to create aws session")
	}

	opts := &token.GetTokenOptions{
		Region:    configFlags.CommonFlags.Region,
		ClusterID: configFlags.ClusterName,
		Session:   sess,
	}
	tok, err := gen.GetWithOptions(opts)
	if err != nil {
		return errors.Wrap(err, "error generating token")
	}

	tokenOutput := &tokenOutput{
		Kind:       "ExecCredential",
		ApiVersion: "client.authentication.k8s.io/v1",
		Status: struct {
			ExpirationTimestamp string `json:"expirationTimestamp"`
			Token               string `json:"token"`
		}{
			ExpirationTimestamp: tok.Expiration.Format("2006-01-02T15:04:05Z"),
			Token:               tok.Token,
		},
	}
	jsonOutput, err := json.Marshal(tokenOutput)
	if err != nil {
		return errors.Wrap(err, "error marshaling json")
	}

	fmt.Println(string(jsonOutput))
	return nil
}
