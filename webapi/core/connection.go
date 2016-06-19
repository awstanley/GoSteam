// Copyright 2016 A.W. Stanley All rights reserved.
// Use of this source code is governed by a BSD-style
// licence that can be found in the LICENCE.md file.

package core

import (
	"fmt"
	"net/http"
	"time"
)

// A Connection object represents a connection to the WebAPI
type Connection struct {
	// API Key
	key string

	// If true, the partner API is used (and secure is forced)
	partner bool

	// If true, HTTPS is used.
	secure bool

	// HTTP Client (for a level of control)
	client *http.Client

	// A buffered baseURI
	baseURI string
}

// IsPartner returns true if the connection is a partner connection.
func (conn *Connection) IsPartner() bool {
	return conn.partner
}

// IsSecure returns true if the connection is an HTTPS connection.
func (conn *Connection) IsSecure() bool {
	return conn.secure
}

// HasKey returns true if an API key is stored.
// Validity is not ensured.
func (conn *Connection) HasKey() bool {
	return len(conn.key) == 32
}

// NewConnection creates a new connection to the Steam API and builds the internal capabilities.
//
// The key parameter is used to specify the user's API key.
// The useSecureProtocol parameter specifies whether or not the API should use HTTPS or HTTP.
// The partner parameter specifies whether or not the API should use the partner endpoint (usually not).
// The fetch parameter specifies whether or not the API should be fetched, if not the user is
// expected to load a copy.
func NewConnection(key string, useSecureProtocol bool, partner bool) (conn *Connection) {

	// First, assume it's valid
	conn = &Connection{
		key:     key,
		partner: partner,
		secure:  useSecureProtocol,
		client: &http.Client{
			CheckRedirect: nil,             // not required at present
			Jar:           nil,             // not required at present
			Timeout:       5 * time.Second, // 5 second timeout
		},
	}

	// Build the baseURI
	base := publicAPI
	if partner {
		// Partner queries must be secure
		conn.secure = true
		base = partnerAPI
	}
	proto := "http://"
	if conn.secure {
		proto = "https://"
	}
	conn.baseURI = fmt.Sprintf("%s%s/", proto, base)

	return conn
}
