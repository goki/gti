// Copyright (c) 2023, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gtigen

import (
	"os"
	"strings"
	"testing"
	"text/template"

	"goki.dev/grease"
	"goki.dev/gti"
	"goki.dev/gti/gtigen/testdata"
	"goki.dev/ordmap"
)

func TestGenerate(t *testing.T) {
	c := &Config{
		AddMethods: true,
		AddFuncs:   true,
		InterfaceConfigs: ordmap.Make([]ordmap.KeyVal[string, *Config]{{
			"fmt.Stringer", &Config{
				AddTypes: true,
				TypeVar:  true,
				Instance: true,
				Setters:  true,
				Templates: []*template.Template{
					template.Must(template.New("Stringer").Parse(`
					func (t *{{.Name}}) MyCustomFuncForStringers(a any) error {
						return nil
					}`)),
				},
			},
		}}),
	}
	err := grease.SetFromDefaults(c)
	if err != nil {
		t.Errorf("programmer error: error setting config from default tags: %v", err)
	}
	c.Dir = "./testdata"
	err = Generate(c)
	if err != nil {
		t.Errorf("error while generating: %v", err)
	}
	have, err := os.ReadFile("testdata/gtigen.go")
	if err != nil {
		t.Errorf("error while reading generated file: %v", err)
	}
	want, err := os.ReadFile("testdata/gtigen.golden")
	if err != nil {
		t.Errorf("error while reading golden file: %v", err)
	}
	// ignore first line, which has "Code generated by" message
	// that can change based on where go test is ran.
	_, shave, got := strings.Cut(string(have), "\n")
	if !got {
		t.Errorf("expected string with newline in testdata/gtigen.go, but got %q", have)
	}
	_, swant, got := strings.Cut(string(want), "\n")
	if !got {
		t.Errorf("expected string with newline in testdata/gtigen.golden, but got %q", want)
	}
	if shave != swant {
		t.Errorf("expected generated file and expected file to be the same after the first line, but they are not (compare ./testdata/gtigen.go and ./testdata/gtigen.golden to see the difference)")
	}
}

func TestPerson(t *testing.T) {
	want := testdata.PersonType
	have := gti.TypeByName("goki.dev/gti/gtigen/testdata.Person")
	if have != want {
		t.Errorf("expected TypeByName to return %v, but got %v", want, have)
	}
	have = gti.TypeByValue(testdata.Person{})
	if have != want {
		t.Errorf("expected TypeByValue to return %v, but got %v", want, have)
	}
	if _, ok := have.Instance.(*testdata.Person); !ok {
		t.Errorf("expected instance to be a Person, but it is a %T (value %v)", have.Instance, have.Instance)
	}
	if have.Name != "goki.dev/gti/gtigen/testdata.Person" {
		t.Errorf("expected name to be 'goki.dev/gti/gtigen/testdata.Person', but got %s", have.Name)
	}
	if len(have.Directives) != 2 {
		t.Errorf("expected 1 directive, but got %d", len(have.Directives))
	}
	if have.Fields.Len() != 3 {
		t.Errorf("expected 2 fields, but got %d", have.Fields.Len())
	}
	if have.Embeds.Len() != 1 {
		t.Errorf("expected 0 embeds, but got %v", have.Embeds.Len())
	}
	if have.Methods.Len() != 1 {
		t.Errorf("expected 1 method, but got %d", have.Methods.Len())
	}

}
