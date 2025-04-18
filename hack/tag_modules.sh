#! /bin/bash
set -euo pipefail

if [[ "$#" -ne 2 ]] ; then
  echo "Usage: ${0##*/} {remote} {tag}"
  echo "This command tags all modules as {module_name}/{tag} then pushes these tags to the selected remote."
  echo "This is used whenever we want to level all modules at the same version."
  exit 1
fi

remote="$1"
tag="$2"
root="$(git rev-parse --show-toplevel)"
declare -a all_tags

cd "${root}"
echo "Tagging all modules in repo ${root##*/}..."

while read module_location ; do
  module_dir=${module_location%/*}
  base_tag="${module_dir#*/}"
  if [[ "${base_tag}" ==  "." ]] ; then
    module_tag="${tag}" # e.g. "v0.24.0"
  else
    module_tag="${base_tag}/${tag}" # e.g. "mangling/v0.24.0"
  fi
  all_tags+=("${module_tag}")
  git tag "${module_tag}"
done < <(find . -name \*.mod | grep -v "\.git" | sort | uniq)

echo "Pushing tags to ${remote}: ${all_tags[@]}"
git push "${remote}" ${all_tags[@]}
