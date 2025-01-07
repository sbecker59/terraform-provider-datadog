#!/usr/bin/env bash
# Generate the documentation using tfplugindocs and remove changes to files that shouldn't change

# Add here the files to be excluded from the doc generation
exclude_files=(
  "docs/resources/integration_aws_account.md"
  "docs/resources/role.md"
)

# Check if manual changes were made to any excluded files and exit
# otherwise these will be lost with `tfplugindocs`
if [ "${#exclude_files[@]}" -ne 0 ] && [ "$(git status --porcelain "${exclude_files[@]}")" ]; then
  echo "Uncommitted changes were detected to the following files. These aren't autogenerated, please commit or stash these changes and try again"
  echo $(git status --porcelain "${exclude_files[@]}")
  exit 1
fi

tfplugindocs

# Remove the changes to files we don't autogenerate
git checkout HEAD -- "${exclude_files[@]}"
