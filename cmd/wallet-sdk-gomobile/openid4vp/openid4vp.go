/*
Copyright Avast Software. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// Package openid4vp contains functionality for doing OpenID4VP operations.
package openid4vp

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/aries-framework-go/pkg/doc/jwt"
	"github.com/hyperledger/aries-framework-go/pkg/doc/presexch"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/piprate/json-gold/ld"

	"github.com/trustbloc/wallet-sdk/cmd/wallet-sdk-gomobile/api"
	"github.com/trustbloc/wallet-sdk/cmd/wallet-sdk-gomobile/credential"
	"github.com/trustbloc/wallet-sdk/cmd/wallet-sdk-gomobile/walleterror"
	"github.com/trustbloc/wallet-sdk/cmd/wallet-sdk-gomobile/wrapper"
	"github.com/trustbloc/wallet-sdk/pkg/common"
	"github.com/trustbloc/wallet-sdk/pkg/openid4vp"
)

type goAPIOpenID4VP interface {
	GetQuery() (*presexch.PresentationDefinition, error)
	PresentCredential(credentials []*verifiable.Credential) error
}

// Interaction represents a single OpenID4VP interaction between a wallet and a verifier. The methods defined on this
// object are used to help guide the calling code through the OpenID4VP flow.
type Interaction struct {
	keyHandleReader  api.KeyReader
	crypto           api.Crypto
	ldDocumentLoader api.LDDocumentLoader
	goAPIOpenID4VP   goAPIOpenID4VP
	didResolver      api.DIDResolver
	inquirer         *credential.Inquirer
}

// ClientConfig contains various parameters for an OpenID4VP Interaction.
// ActivityLogger is optional, but if provided then activities will be logged there.
// If not provided, then no activities will be logged.
type ClientConfig struct {
	KeyHandleReader                  api.KeyReader
	Crypto                           api.Crypto
	DIDRes                           api.DIDResolver
	DocumentLoader                   api.LDDocumentLoader
	ActivityLogger                   api.ActivityLogger
	MetricsLogger                    api.MetricsLogger
	additionalHeaders                api.Headers
	disableHTTPClientTLSVerification bool
}

// NewClientConfig creates the client config object.
// ActivityLogger is optional, but if provided then activities will be logged there.
// If ldDocumentLoader is set to nil, then a network-based loader will be used.
// If not provided, then no activities will be logged.
func NewClientConfig(keyHandleReader api.KeyReader, crypto api.Crypto,
	didResolver api.DIDResolver, ldDocumentLoader api.LDDocumentLoader, activityLogger api.ActivityLogger,
) *ClientConfig {
	return &ClientConfig{
		KeyHandleReader: keyHandleReader,
		Crypto:          crypto,
		DIDRes:          didResolver,
		DocumentLoader:  ldDocumentLoader,
		ActivityLogger:  activityLogger,
	}
}

// AddHeaders adds the given HTTP headers to all REST calls made to the verifier during the OpenID4VP flow.
func (c *ClientConfig) AddHeaders(headers *api.Headers) {
	headersAsArray := headers.GetAll()

	for i := range headersAsArray {
		c.additionalHeaders.Add(&headersAsArray[i])
	}
}

// DisableHTTPClientTLSVerify disables tls verification, should be used only for test purposes.
func (c *ClientConfig) DisableHTTPClientTLSVerify() {
	c.disableHTTPClientTLSVerification = true
}

// NewInteraction creates a new OpenID4VP Interaction.
// The methods defined on this object are used to help guide the calling code through the OpenID4CI flow.
// If activityLogger is nil, then no activity logging will take place.
func NewInteraction(authorizationRequest string, config *ClientConfig) *Interaction {
	jwtVerifier := jwt.NewVerifier(jwt.KeyResolverFunc(
		common.NewVDRKeyResolver(&wrapper.VDRResolverWrapper{
			DIDResolver: config.DIDRes,
		}).PublicKeyFetcher()))

	httpClient := wrapper.NewHTTPClient()
	httpClient.AddHeaders(&config.additionalHeaders)
	httpClient.DisableTLSVerification = config.disableHTTPClientTLSVerification

	opts := []openid4vp.Opt{openid4vp.WithHTTPClient(httpClient)}

	if config.ActivityLogger != nil {
		mobileActivityLoggerWrapper := &wrapper.MobileActivityLoggerWrapper{MobileAPIActivityLogger: config.ActivityLogger}

		opts = append(opts, openid4vp.WithActivityLogger(mobileActivityLoggerWrapper))
	}

	if config.MetricsLogger != nil {
		mobileMetricsLoggerWrapper := &wrapper.MobileMetricsLoggerWrapper{MobileAPIMetricsLogger: config.MetricsLogger}

		opts = append(opts, openid4vp.WithMetricsLogger(mobileMetricsLoggerWrapper))
	}

	var goAPIDocumentLoader ld.DocumentLoader

	if config.DocumentLoader != nil {
		goAPIDocumentLoader = &wrapper.DocumentLoaderWrapper{
			DocumentLoader: config.DocumentLoader,
		}
	}

	inquirer := credential.NewInquirer(config.DocumentLoader)

	return &Interaction{
		keyHandleReader:  config.KeyHandleReader,
		ldDocumentLoader: config.DocumentLoader,
		crypto:           config.Crypto,
		goAPIOpenID4VP: openid4vp.New(
			authorizationRequest,
			jwtVerifier,
			&wrapper.VDRResolverWrapper{DIDResolver: config.DIDRes},
			config.Crypto,
			goAPIDocumentLoader,
			opts...,
		),
		didResolver: config.DIDRes,
		inquirer:    inquirer,
	}
}

// GetQuery creates query based on authorization request data.
func (o *Interaction) GetQuery() ([]byte, error) {
	presentationDefinition, err := o.goAPIOpenID4VP.GetQuery()
	if err != nil {
		return nil, walleterror.ToMobileError(err)
	}

	pdBytes, err := json.Marshal(presentationDefinition)
	if err != nil {
		return nil, walleterror.ToMobileError(
			fmt.Errorf("presentation definition marshal: %w", err))
	}

	return pdBytes, nil
}

// PresentCredential presents credentials to redirect uri from request object.
func (o *Interaction) PresentCredential(credentials *api.VerifiableCredentialsArray) error {
	return walleterror.ToMobileError(o.goAPIOpenID4VP.PresentCredential(unwrapVCs(credentials)))
}

func unwrapVCs(vcs *api.VerifiableCredentialsArray) []*verifiable.Credential {
	var credentials []*verifiable.Credential

	for i := 0; i < vcs.Length(); i++ {
		credentials = append(credentials, vcs.AtIndex(i).VC)
	}

	return credentials
}
