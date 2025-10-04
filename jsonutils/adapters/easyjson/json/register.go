// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package json

import (
	"fmt"
	"reflect"

	"github.com/go-openapi/swag/jsonutils/adapters/ifaces"
	"github.com/mailru/easyjson"
)

// Register the easyjson implementation of a [ifaces.Adapter] to the an [ifaces.Registrar],
// e.g. the global registry [github.com/go-openapi/swag/jsonutils/adapters.Registry].
//
// [Register] calls [ifaces.Registrar.RegisterFor].
//
// Some optional features proposed by the [jwriter.Writer] and [jlexer.Lexer] are available. See [Option].
func Register(dispatcher ifaces.Registrar, opts ...Option) {
	t := reflect.TypeOf(Adapter{})
	var o options
	for _, apply := range opts {
		apply(&o)
	}

	dispatcher.RegisterFor(
		ifaces.RegistryEntry{
			Who:         fmt.Sprintf("%s.%s", t.PkgPath(), t.Name()),
			What:        ifaces.AllCapabilities,
			Constructor: BorrowAdapterIface,
			Support:     support,
		})
}

func support(capability ifaces.Capability, value any) bool {
	switch capability {
	case ifaces.CapabilityMarshalJSON, ifaces.CapabilityOrderedMarshalJSON:
		_, ok := value.(easyjson.Marshaler)
		return ok
	case ifaces.CapabilityUnmarshalJSON, ifaces.CapabilityOrderedUnmarshalJSON:
		_, ok := value.(easyjson.Unmarshaler)
		return ok
	case ifaces.CapabilityOrderedMap:
		return true
	default:
		return false
	}
}
