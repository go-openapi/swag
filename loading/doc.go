// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

// Package loading provides tools to load a file from http or from a local file system.
//
// # Security
//
// By default, the local loader reads any path the process can access, including absolute
// paths and "file://" URIs (for example "file:///etc/passwd"). Applications that pass
// untrusted input to [LoadFromFileOrHTTP], [JSONDoc] (or to downstream consumers such as
// go-openapi/loads) must confine local loading to a trusted directory.
//
// Use [WithRoot] to do so: it resolves every requested path relative to a chosen directory
// and rejects anything that escapes it, including via symlink. It is built on [os.Root]
// and is therefore safer than passing an [os.DirFS] to [WithFS], which does not block
// symlink escapes.
package loading
