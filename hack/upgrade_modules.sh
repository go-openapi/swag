#! /bin/bash
new_tag=$1

cd "$(git rev-parse --show-toplevel)"
while read -r dir
do
  pushd $dir
  go list -deps -test \
    -f '{{ if .DepOnly }}{{ with .Module }}{{ .Path }}{{ end }}{{ end }}'|\
  sort -u|\
  grep "go-openapi/swag"|\
  while read -r module ; do
    go mod edit -require "${module}@${new_tag}"
  done

  go mod tidy
  popd
done < <(find . -name go.mod |xargs dirname)

go work sync
