# Roadmap

A few ideas on the todo list:

* [x] Complete the split of dependencies to isolate easyjson from the rest
* [x] Improve CI to reduce needed tests
* [x] Replace dependency to `gopkg.in/yaml.v3` (`yamlutil`)
* [ ] Improve mangling utilities (improve readability, support for capitalized words,
      better word substitution for non-letter symbols...)
* [ ] Move back to this common shared pot a few of the technical features introduced by go-swagger independently
      (e.g. mangle go package names, search package with go modules support, ...)
* [ ] Apply a similar mono-repo approach to `go-openapi/strfmt` which suffer from similar woes: bloated API,
      imposed dependency to some database driver.
* [ ] Adapt `go-swagger` (incl. generated code) to the new `swag` API.
* [ ] Factorize some tests, as there is a lot of redundant testing code in `jsonutils`
* [ ] Benchmark & profiling: publish independently the tool built to analyze and chart benchmarks (e.g. similar to `benchvisual`)
* [ ] more thorough testing for nil / null case
* [ ] ci pipeline to manage releases
* [ ] cleaner mockery generation (doesn't work out of the box for all sub-modules)
