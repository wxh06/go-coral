// Copyright (c) 2026 Xinhe Wang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package contract is the seam between the Coral client, github.com/wxh06/go-coral, and whatever
// seals its requests.
//
// Coral, the service behind the Nintendo Switch Online app, accepts only encrypted request
// bodies, and its authentication endpoints also carry a signature, f, that is bound to the sealed
// request: the two cannot be produced separately. The swappable unit is therefore the whole
// envelope. A [Provider] seals a [Request] into the bytes that go on the wire and opens the reply;
// the client knows nothing of the crypto, and this package carries none of it.
//
// It is types only, in a module of its own, so the client and a provider can each depend on it
// without depending on the other. Every change here moves every provider, so it stays minimal.
package contract

import "context"

// HashMethod is the service's own hash_method numeral: which f scheme seals a request, if any.
// The values are the wire contract, not an enumeration order.
type HashMethod int

const (
	// HashMethodNone seals an ordinary endpoint: no f, the body goes through as is.
	HashMethodNone HashMethod = iota
	// HashMethodAuth is Account/Login and Account/GetToken; the token is the Nintendo Account
	// id_token.
	HashMethodAuth
	// HashMethodWebService is Game/GetWebServiceToken; the token is the Coral access token.
	HashMethodWebService
)

// Request is what a Provider seals: the superset every provider reads from, so a caller builds
// the same Request whichever provider it holds. A provider ignores the fields it has no use for.
type Request struct {
	// URL is the absolute Coral URL the envelope will be posted to. A provider that seals
	// remotely sends it along; one that seals in-process never reads it.
	URL string

	// Token is what f is computed from under HashMethodAuth (the Nintendo Account id_token) and
	// HashMethodWebService (the Coral access token). Under HashMethodNone it is the Bearer the
	// request will be sent with, or empty when the request sends none.
	Token string

	// Body is the plain JSON body, an object. Under HashMethodAuth and HashMethodWebService its
	// "parameter" object carries a requestId (a lowercase UUID v4) and a timestamp (Unix
	// milliseconds) that the caller has filled in; the provider inserts f beside them, replacing
	// both if its f scheme generates its own. A provider may re-encode the body, and the server
	// reads it as an unordered object, so nothing may depend on its byte layout.
	Body []byte

	// HashMethod picks the f scheme, if any.
	HashMethod HashMethod

	// NAID is the Nintendo Account id, the "id" of https://api.accounts.nintendo.com/2.0.0/users/me.
	// A provider that seals remotely sends it for validation: required under
	// HashMethodWebService, recommended under HashMethodAuth, unused under HashMethodNone.
	NAID string

	// CoralUserID is the Coral user id, Account/Login's result.user.id, in decimal. Only
	// HashMethodWebService uses it, so it is empty until the first Login has returned one. A
	// provider that seals remotely sends it for validation; it is recommended, not required.
	CoralUserID string
}

// Provider seals Coral requests and opens their replies.
//
// EncryptRequest returns the envelope that is posted as the request body. DecryptResponse opens
// a reply's body into the JSON it carries; a reply is self-contained, so nothing from the request
// is needed. AppVersion is the Nintendo Switch Online app version the envelopes claim; the client
// sends the same string in its own headers, so it has to come from here.
//
// An implementation must be safe for concurrent use: one Provider serves every request of a
// client.
type Provider interface {
	EncryptRequest(ctx context.Context, req Request) (envelope []byte, err error)
	DecryptResponse(ctx context.Context, envelope []byte) (body []byte, err error)
	AppVersion() string
}
