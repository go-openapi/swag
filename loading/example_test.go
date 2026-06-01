// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package loading_test

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"syscall"

	"github.com/go-openapi/swag/loading"
)

// errForbiddenAddr is returned by the dial guard when a destination is not allowed.
var errForbiddenAddr = errors.New("blocked dial to a forbidden address")

// ExampleLoadFromFileOrHTTP_restrictNetwork shows how to confine remote spec loading so a
// caller-controlled URL cannot reach loopback, private, or link-local (cloud metadata)
// addresses.
//
// The [net.Dialer] Control hook runs after DNS resolution and before connect, on every
// connection, so the check also covers HTTP redirects and DNS rebinding — neither of which
// a URL-string allowlist can defend against. Here a loopback test server stands in for an
// internal endpoint that the guard must refuse to reach.
func ExampleLoadFromFileOrHTTP_restrictNetwork() {
	control := func(_, address string, _ syscall.RawConn) error {
		host, _, err := net.SplitHostPort(address)
		if err != nil {
			return err
		}
		addr, err := netip.ParseAddr(host)
		if err != nil {
			return err
		}
		if a := addr.Unmap(); a.IsLoopback() || a.IsPrivate() || a.IsLinkLocalUnicast() || a.IsUnspecified() {
			return errForbiddenAddr
		}

		return nil
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{Control: control}).DialContext,
		},
	}

	// An internal service the application must not let untrusted input reach.
	internal := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("internal secret"))
	}))
	defer internal.Close()

	// internal.URL is a loopback address (the untrusted URL in a real attack).
	_, err := loading.LoadFromFileOrHTTP(internal.URL, loading.WithHTTPClient(client))
	fmt.Println("blocked:", errors.Is(err, errForbiddenAddr))

	// Output:
	// blocked: true
}
