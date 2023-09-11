// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gtigen

// Config contains the configuration information
// used by gtigen
type Config struct {

	// [def: .] the source directory to run enumgen on (can be set to multiple through paths like ./...)
	Dir string `def:"." desc:"the source directory to run enumgen on (can be set to multiple through paths like ./...)"`

	// [def: gtigen.go] the output file location relative to the package on which enumgen is being called
	Output string `def:"gtigen.go" desc:"the output file location relative to the package on which enumgen is being called"`

	TypeReg  bool
	FuncReg  bool
	VarReg   bool
	ConstReg bool

	// whether to generate an instance of the type(s)
	Instances bool `desc:"whether to generate an instance of the type(s)"`

	// whether to generate a global type variable of the form 'TypeNameType'
	TypeVar bool `desc:"whether to generate a global type variable of the form 'TypeNameType'"`
}