// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

// Package easyjson exposes a JSON adapter
// that leverages the [easyjson] serializer library.
//
// It ships as an independent go module.
//
// This library is significantly faster than the standard
// library, provided the data types implement its specific
// interfaces [easyjson.Marshaler] and [easyjson.Unmarshaler].
package easyjson

import (
	_ "github.com/mailru/easyjson" // for documentation purpose only
)
