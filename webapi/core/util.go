// Copyright 2016 A.W. Stanley All rights reserved.
// Use of this source code is governed by a BSD-style
// licence that can be found in the LICENCE.md file.

package core

import (
	"fmt"
	"net/url"
)

// Parameters is a trivial wrapper around url.Values
type Parameters struct {
	Values url.Values
}

// NewParameters creates a new parameter set.
func NewParameters() *Parameters {
	return &Parameters{
		Values: url.Values{},
	}
}

// Encode encodes the Parameters object with URL encoding
func (p *Parameters) Encode() string {
	return p.Values.Encode()
}

// SetKey sets the API key
func (p *Parameters) SetKey(key string) {
	p.Values.Set("key", key)
}

// AddString adds a string.
func (p *Parameters) AddString(name string, value string) {
	p.Values.Add(name, value)
}

// AddBytes stores bytes.
func (p *Parameters) AddBytes(name string, value []byte) {
	p.Values.Add(name, string(value))
}

// AddInt32 adds an int32
func (p *Parameters) AddInt32(name string, value int32) {
	p.Values.Add(name, fmt.Sprintf("%d", value))
}

// AddUInt32 adds a uint32
func (p *Parameters) AddUInt32(name string, value uint32) {
	p.Values.Add(name, fmt.Sprintf("%d", value))
}

// AddUInt64 adds a uint64
func (p *Parameters) AddUInt64(name string, value uint64) {
	p.Values.Add(name, fmt.Sprintf("%d", value))
}

// AddFloat32 adds a float32
func (p *Parameters) AddFloat32(name string, value float32) {
	p.Values.Add(name, fmt.Sprintf("%g", value))
}

// AddBoolean adds a boolean
func (p *Parameters) AddBoolean(name string, value bool) {
	if value {
		p.Values.Add(name, "true")
	} else {
		p.Values.Add(name, "false")
	}
}
