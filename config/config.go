//go:build linux

// Copyright (c) 2024 Generic API Server All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package config 설정 패키지
*/
package config

const (
	ModuleName = "openkms"
)

const (
	PidFilePath = "var/" + ModuleName + ".pid"
)
