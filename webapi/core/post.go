// Copyright 2016 A.W. Stanley All rights reserved.
// Use of this source code is governed by a BSD-style
// licence that can be found in the LICENCE.md file.

package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Post performs a POST request against the base using the stored key (if required)
func (conn *Connection) Post(uri string, params *Parameters, requireKey bool) (content []byte, err error) {

	if conn.partner {
		requireKey = true
	}

	if requireKey {
		params.SetKey(conn.key)
	}
	uri = fmt.Sprintf("%s%s", conn.baseURI, uri)

	payload := params.Encode()

	request, _ := http.NewRequest("POST", uri, bytes.NewBufferString(payload))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	response, err := conn.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err = ioutil.ReadAll(response.Body)

	// Return any error that arose from the body
	return content, err
}
