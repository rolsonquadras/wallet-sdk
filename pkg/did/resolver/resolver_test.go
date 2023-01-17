/*
Copyright Avast Software. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package resolver_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hyperledger/aries-framework-go-ext/component/vdr/jwk"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/stretchr/testify/require"

	"github.com/trustbloc/wallet-sdk/pkg/did/resolver"
)

const (
	docID = "did:ion:test"
	doc   = `{
  "@context": ["https://w3id.org/did/v1","https://w3id.org/did/v2"],
  "id": "did:ion:test"
}`
)

func TestDIDResolver(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		didResolver, err := resolver.NewDIDResolver("")
		require.NoError(t, err)

		didDocResolution, err := didResolver.Resolve("did:key:z6MkjfbzWitsSUyFMTbBUSWNsJBHR7BefFp1WmABE3kRw8Qr")
		require.NoError(t, err)
		require.NotEmpty(t, didDocResolution)
	})

	t.Run("did jwk", func(t *testing.T) {
		vdr := jwk.New()
		didDoc, err := vdr.Create(&did.Doc{VerificationMethod: []did.VerificationMethod{},
			AssertionMethod: []did.Verification{{
				Relationship: did.AssertionMethod,
				Embedded:     true,
			}}})
		require.NoError(t, err)

		didResolver, err := resolver.NewDIDResolver("")
		require.NoError(t, err)

		didDocResolution, err := didResolver.Resolve(didDoc.DIDDocument.ID)
		require.NoError(t, err)
		require.NotEmpty(t, didDocResolution)
	})

	t.Run("httpbinding initialization error", func(t *testing.T) {
		didResolver, err := resolver.NewDIDResolver("not a uri")
		require.Error(t, err)
		require.Nil(t, didResolver)
		require.Contains(t, err.Error(), "failed to initialize client for DID resolution server")
	})

	t.Run("httpbinding resolve", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Header().Add("Content-type", "application/did+ld+json")
			res.WriteHeader(http.StatusOK)

			_, err := res.Write([]byte(doc))
			require.NoError(t, err)
		}))

		defer func() { testServer.Close() }()

		didResolver, err := resolver.NewDIDResolver(testServer.URL)
		require.NoError(t, err)

		didDocResolution, err := didResolver.Resolve(docID)
		require.NoError(t, err)
		require.NotNil(t, didDocResolution)
		require.NotNil(t, didDocResolution.DIDDocument)
		require.Equal(t, docID, didDocResolution.DIDDocument.ID)
	})
}

func TestDIDResolver_InvalidDID(t *testing.T) {
	didResolver, err := resolver.NewDIDResolver("")
	require.NoError(t, err)

	didDocResolution, err := didResolver.Resolve("did:example:abc")
	require.Error(t, err)
	require.EqualError(t, err, "resolve did:example:abc : did method example not supported for vdr")
	require.Empty(t, didDocResolution)
}
