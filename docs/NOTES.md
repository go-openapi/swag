# pre-v0.26.0 release notes

## v0.25.4

**mangling**

Bug fix

* [x] mangler may panic with pluralized overlapping initialisms

Tests

* [x] introduced fuzz tests

## v0.25.3

**mangling**

Bug fix

* [x] mangler may panic with pluralized initialisms

## v0.25.2

Minor changes due to internal maintenance that don't affect the behavior of the library.

* [x] removed indirect test dependencies by switching all tests to `go-openapi/testify`,
  a fork of `stretch/testify` with zero-dependencies.
* [x] improvements to CI to catch test reports.
* [x] modernized licensing annotations in source code, using the more compact SPDX annotations
  rather than the full license terms.
* [x] simplified a bit JSON & YAML testing by using newly available assertions
* started the journey to an OpenSSF score card badge:
  * [x] explicated permissions in CI workflows
  * [x] published security policy
  * pinned dependencies to github actions
  * introduced fuzzing in tests

## v0.25.1

* fixes a data race that could occur when using the standard library implementation of a JSON ordered map

## v0.25.0

**New with this release**:

* requires `go1.24`, as iterators are being introduced
* removes the dependency to `mailru/easyjson` by default (#68)
  * functionality remains the same, but performance may somewhat degrade for applications
    that relied on `easyjson`
  * users of the JSON or YAML utilities who want to use `easyjson` as their preferred JSON serializer library
    will be able to do so by registering this the corresponding JSON adapter at runtime. See below.
  * ordered keys in JSON and YAML objects: this feature used to rely solely on `easyjson`.
    With this release, an implementation relying on the standard `encoding/json` is provided.
  * an independent [benchmark](../jsonutils/adapters/testintegration/benchmarks/README.md) to compare the different adapters
* improves the "float is integer" check (`conv.IsFloat64AJSONInteger`) (#59)
* removes the _direct_ dependency to `gopkg.in/yaml.v3` (indirect dependency is still incurred through `stretchr/testify`) (#127)
* exposed `conv.IsNil()` (previously kept private): a safe nil check (accounting for the "non-nil interface with nil value" nonsensical go trick)

## v0.24.0

With this release, we have largely modernized the API of `swag`:

* The traditional `swag` API is still supported: code that imports `swag` will still
  compile and work the same.
* A deprecation notice is published to encourage consumers of this library to adopt
  the newer API
* **Deprecation notice**
  * configuration through global variables is now deprecated, in favor of options passed as parameters
  * all helper functions are moved to more specialized packages, which are exposed as
    go modules. Importing such a module would reduce the footprint of dependencies.
  * _all_ functions, variables, constants exposed by the deprecated API have now moved, so
    that consumers of the new API no longer need to import github.com/go-openapi/swag, but
    should import the desired sub-module(s).

**New with this release**:

* [x] type converters and pointer to value helpers now support generic types
* [x] name mangling now support pluralized initialisms (issue #46)
      Strings like "contact IDs" are now recognized as such a plural form and mangled as a linter would expect.
* [x] performance: small improvements to reduce the overhead of convert/format wrappers (see issues #110, or PR #108)
* [x] performance: name mangling utilities run ~ 10% faster (PR #115)
