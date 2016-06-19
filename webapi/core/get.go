// Copyright 2016 A.W. Stanley All rights reserved.
// Use of this source code is governed by a BSD-style
// licence that can be found in the LICENCE.md file.

package core

import (
	"fmt"
	"io/ioutil"
)

// Get performs a Get request against the base using the stored key (if required)
func (conn *Connection) Get(uri string, params *Parameters, requireKey bool) (content []byte, err error) {

	if conn.partner {
		requireKey = true
	}

	if requireKey {
		params.SetKey(conn.key)
	}
	uri = fmt.Sprintf("%s%s?%s", conn.baseURI, uri, params.Encode())

	response, err := conn.client.Get(uri)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err = ioutil.ReadAll(response.Body)

	// Return any error that arose from the body
	return content, err
}
