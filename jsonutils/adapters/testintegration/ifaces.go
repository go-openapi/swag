// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package testintegration

import (
	"github.com/mailru/easyjson"
)

// EJMarshaler wraps [easyjson.Marshaler]
type EJMarshaler interface {
	easyjson.Marshaler
}

// EJUnmarshaler wraps [easyjson.Unmarshaler]
type EJUnmarshaler interface {
	easyjson.Unmarshaler
}
