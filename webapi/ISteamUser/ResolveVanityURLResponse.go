// Copyright 2016 A.W. Stanley All rights reserved.
// Use of this source code is governed by a BSD-style
// licence that can be found in the LICENCE.md file.

package ISteamUser

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ResolveVanityURLV1Response respresents the JSON return value.
type ResolveVanityURLV1Response struct {
	SteamID uint64
	Success int
}

// Decode transforms the raw byte content into a neat struct.
func (res *ResolveVanityURLV1Response) Decode(contents []byte) error {
	type SuccessInner struct {
		SteamID string `json:"steamid"`
		Success int    `json:"success"`
	}

	type SuccessResponse struct {
		InnerStruct SuccessInner `json:"response"`
	}

	type FailInner struct {
		Success int    `json:"success"`
		Message string `json:"message"`
	}

	type FailResponse struct {
		InnerStruct FailInner `json:"response"`
	}

	response := SuccessResponse{}

	err := json.Unmarshal(contents, &response)

	// Catch the failure state before it gets out of hand.
	if err != nil || response.InnerStruct.Success != 1 {
		failure := FailResponse{}
		err = json.Unmarshal(contents, &failure)
		if err != nil {
			return err // this is *bad*
		}

		// Failure fallback
		return fmt.Errorf("query returned with non-success value (%d) and message '%s'\n",
			failure.InnerStruct.Success,
			failure.InnerStruct.Message,
		)
	}

	res.SteamID, err = strconv.ParseUint(response.InnerStruct.SteamID, 10, 64)
	if err != nil {
		return err
	}
	res.Success = response.InnerStruct.Success
	return nil
}
