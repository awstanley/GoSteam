// Copyright 2016 A.W. Stanley All rights reserved.
// Use of this source code is governed by a BSD-style
// licence that can be found in the LICENCE.md file.

package main

// This trivial example performs a call to:
// https://api.steampowered.com/ISteamUser/ResolveVanityURL/v1/

import (
	"encoding/json"                                  // json support
	"flag"                                           // args
	"fmt"                                            // format support
	"github.com/awstanley/GoSteam/webapi/ISteamUser" // ISteamUser
	"github.com/awstanley/GoSteam/webapi/core"       // the core
	"os"                                             // to get the binary name in Usage
	"strconv"                                        // strconv helps manipulate non-strings into strings
)

// The type of vanity URL. 1 (default): Individual profile, 2: Group, 3: Official game group

func Usage() {
	fmt.Println("Usage:")
	fmt.Printf("%s --key=<key>\n\n", os.Args[0])
	fmt.Println("  --vanity=<vanity>")
}

func main() {

	// API key
	apiKey := flag.String("key", "", "Steam WebAPI key")

	// e.g. my url http://steamcommunity.com/id/swixel/ = "swixel"
	// This should return 76561197993978679 (unless you change the value)
	vanity := flag.String("vanity", "swixel", "The vanity URL for which you want the SteamID")

	// There's an optional value, but we'll ignore it.
	// urlType := flag.String("url_type")

	flag.Usage = Usage

	flag.Parse()

	if *apiKey == "" {
		fmt.Println("No key provided.")
		return
	}

	if *vanity == "" {
		fmt.Println("No vanity url provided.")
		return
	}

	// Create a connection
	conn := core.NewConnection(*apiKey, true, false)
	if conn == nil {
		fmt.Println("Connection failed to create (this shouldn't happen)")
		return
	}

	// Create the resolver object.
	// Note: normally you'd inline define, but this is to make a point
	resolver := ISteamUser.ResolveVanityURLV1{}
	resolver.Vanityurl = *vanity

	// Call the resolver (passing a connection)
	contents, err := resolver.Call(conn)
	if err != nil {
		fmt.Printf("Resolver failed: %s\n", err)
		return
	}

	// We'll take a look at the data, raw dumping it...
	fmt.Printf("\n---\nRaw JSON (response)\n---\n%s\n\n---\n", string(contents))

	// In this case I actually wrote a handler:
	flag.Usage = Usage
	response := ISteamUser.ResolveVanityURLV1Response{}
	err = response.Decode(contents)
	if err == nil {
		fmt.Printf("[Decode] Received response (success: %d); SteamID: %d\n", response.Success, response.SteamID)
	} else {
		fmt.Printf("[Decode] Failed to decode with error: %s\n", err)
	}

	// But here's how it's normally done..
	// We know the structure from looking at it in the browser,
	// so we can skip iterating blindly over maps of map[string]interface{}.
	// {
	//   "response": {
	//     `success`: 0, // (int, 0 or 1)
	//     `steamid`: "" // (string, value)
	//   }
	// }

	// So we decode into an interface
	var arbitraryDecode interface{}
	err = json.Unmarshal(contents, &arbitraryDecode)
	if err != nil {
		fmt.Printf("[Manual] Failed to decode contents into arbitrary data (shouldn't fail if the above didn't)\n")
		return
	}

	// Cast to a map
	m := arbitraryDecode.(map[string]interface{})

	// Get the inner map ('response'), casting it.
	inner := m["response"].(map[string]interface{})

	// Get "success" as an int (which Go seems to think is float64)
	resultSuccess := inner["success"].(float64)

	if resultSuccess == 1 {
		// Get the string, but since we want it in uint64 format we'll use strconv.
		resultSteamID, err := strconv.ParseUint(inner["steamid"].(string), 10, 64)
		if err != nil {
			// This shouldn't happen unless your result is bad.
			fmt.Printf("[Manual] Failed to parse with strconv: %s\n", err)
			return
		}
		// And we're good.
		fmt.Printf("[Manual] Received response (success: %g); SteamID: %d\n", resultSuccess, resultSteamID)
		return
	}

	// Failure states are fun too ...
	resultMessage := inner["message"].(string)
	fmt.Printf("[Manual] Received non-success code (%g) with message: %s\n", resultSuccess, resultMessage)

}
