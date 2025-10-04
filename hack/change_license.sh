#! /bin/bash
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "${root}"

find . -name \*.go|xargs \
  perl -gi -p -E \
  's/((?:\/\/)|(?:\/\*))\s*Copyright 2015.+limitations\s+under\s+the\s+License\.?(?:\*\/)?/\/\/ SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers\n\/\/ SPDX-License-Identifier: Apache-2.0/sn'

# check
rm -rf /tmp/swap
mkdir -p /tmp/swap
find . -name \*.go |while read file ;do
  if ! $(grep -q SPDX "${file}") ; then
    echo "Missing license in source: ${file} - adding it"
    swapfile="/tmp/swap/$(basename ${file})"
    (cat << EOF
// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

EOF
    ) > "${swapfile}" && \
    cat "${file}">>"${swapfile}"
    mv "${swapfile}" "${file}"
    echo "Done for ${file}"
  fi
done
