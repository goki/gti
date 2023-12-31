// Copyright (c) 2023, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gti

import (
	"reflect"
	"strings"
)

// Type represents a type
type Type struct {
	// Name is the fully package-path-qualified name of the type (eg: goki.dev/gi/v2/gi.Button)
	Name string

	// ShortName is the short, package-qualified name of the type (eg: gi.Button)
	ShortName string

	// IDName is the short, package-unqualified, kebab-case name of the type that is suitable
	// for use in an ID (eg: button)
	IDName string

	// Doc has all of the comment documentation
	// info as one string with directives removed.
	Doc string

	// Directives has the parsed comment directives
	Directives Directives

	// unique type ID number
	ID uint64

	// Methods are available for all types
	Methods *Methods

	// Embedded fields for struct types
	Embeds *Fields

	// Fields for struct types
	Fields *Fields

	// instance of the type
	Instance any

	// All embedded fields (including nested ones) for struct types;
	// not set by gtigen -- HasEmbed automatically compiles it as needed.
	// Key is the ID of the type.
	AllEmbeds map[uint64]*Type
}

func (tp *Type) String() string {
	return tp.Name
}

func (tp *Type) Label() string {
	return tp.ShortName
}

// ReflectType returns the [reflect.Type] for this type, using the Instance
func (tp *Type) ReflectType() reflect.Type {
	if tp.Instance == nil {
		return nil
	}
	return reflect.TypeOf(tp.Instance).Elem()
}

// HasEmbed returns true if this type has the given type
// at any level of embedding depth, including if this type is
// the given type.  The first time called it will Compile
// a map of all embedded types so subsequent calls are very fast.
func (tp *Type) HasEmbed(typ *Type) bool {
	if tp.AllEmbeds == nil {
		tp.CompileEmbeds()
		if tp.AllEmbeds == nil {
			return typ == tp
		}
	}
	if tp == typ {
		return true
	}
	_, has := tp.AllEmbeds[typ.ID]
	return has
}

func (tp *Type) CompileEmbeds() {
	if tp.Embeds == nil {
		return
	}
	rt := tp.ReflectType()
	if rt == nil {
		return
	}
	tp.AllEmbeds = make(map[uint64]*Type)
	for _, em := range tp.Embeds.Order {
		enm := em.Val.Name
		if idx := strings.LastIndex(enm, "."); idx >= 0 {
			enm = enm[idx+1:]
		}
		etf, has := rt.FieldByName(enm)
		if !has {
			continue
		}
		etft := TypeName(etf.Type)
		et := TypeByName(etft)
		if et == nil {
			continue
		}
		tp.AllEmbeds[et.ID] = et
		et.CompileEmbeds()
		if et.AllEmbeds == nil {
			continue
		}
		for id, ct := range et.AllEmbeds {
			tp.AllEmbeds[id] = ct
		}
	}
}
