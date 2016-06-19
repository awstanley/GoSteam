// Copyright 2016 A.W. Stanley All rights reserved.
// Use of this source code is governed by a BSD-style
// licence that can be found in the LICENCE.md file.

package main

// Since the maps represented below are designed to
// be immutable they are stored as instances.

// Root level
type jsonSteamSupportedRoot struct {
	Apilist jsonSteamApiList `json:"apilist"`
}

// Raw JSON apilist
type jsonSteamApiList struct {
	Interfaces []jsonSteamInterface `json:"interfaces"`
}

type jsonSteamInterface struct {
	Name    string            `json:"name"`
	Methods []jsonSteamMethod `json:"methods"`
}

type jsonSteamMethod struct {
	Name       string               `json:"name"`
	Version    int                  `json:"version"`
	Httpmethod string               `json:"httpmethod"`
	Parameters []jsonSteamParameter `json:"parameters"`
}

type jsonSteamParameter struct {
	Name        string `json:"name"`
	VarType     string `json:"type"`
	Optional    bool   `json:"optional"`
	Description string `json:"description,omitempty"`
}

// Storage for the processed API
type apiSteam struct {
	interfaces map[string]*apiSteamInterface
}

func (api *apiSteam) load(root *jsonSteamSupportedRoot) {
	api.interfaces = make(map[string]*apiSteamInterface)

	// Step through
	ifaces := root.Apilist.Interfaces

	// Duplicate interfaces would be fatal, we're not
	// even going to check for it.
	for _, v := range ifaces {
		api.interfaces[v.Name] = &apiSteamInterface{
			methods: make(map[string]*apiSteamMethod),
		}
		api.interfaces[v.Name].load(&v)
	}
}

// Storage for processed interfaces
type apiSteamInterface struct {
	methods map[string]*apiSteamMethod
}

func (api *apiSteamInterface) load(json *jsonSteamInterface) {
	api.methods = make(map[string]*apiSteamMethod)

	// Methods are versioned
	for _, v := range json.Methods {
		method, _ := api.methods[v.Name]
		if method == nil {
			method = &apiSteamMethod{
				make(map[int]*apiSteamVersionedMethod),
			}
		}
		method.load(&v)
		api.methods[v.Name] = method
	}
}

// Storage for processed methods
type apiSteamMethod struct {
	methods map[int]*apiSteamVersionedMethod
}

func (api *apiSteamMethod) load(json *jsonSteamMethod) {
	// No version should exist
	versioned := &apiSteamVersionedMethod{
		params: make(map[string]*apiSteamParameter),
	}
	versioned.load(json)
	versioned.verb = json.Httpmethod
	api.methods[json.Version] = versioned
}

// Stored for versioned methods
type apiSteamVersionedMethod struct {
	verb   string
	params map[string]*apiSteamParameter
}

type apiSteamParameter struct {
	name        string
	varType     string
	optional    bool
	description string
}

func (api *apiSteamVersionedMethod) load(json *jsonSteamMethod) {
	for _, v := range json.Parameters {
		api.params[v.Name] = &apiSteamParameter{
			name:        v.Name,
			varType:     v.VarType,
			optional:    v.Optional,
			description: v.Description,
		}
	}
}
