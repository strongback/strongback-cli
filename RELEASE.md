### Prerequisites

Before running the release, make sure the `master` branch in the [upstream Git repository](https://github.com/strongback/strongback-cli) contains the code to be released, and that your local repository _has no changes and reflects the upstream repository_. Open a terminal and go to the `cli` directory (the top of your local repository).

### Perform the release

First, use the following script to bump the version number:

    $ bin/bump-version.sh <component>

where `<component>` is one of:

* `major` to perform a major release with breaking changes
* `minor` to perform a minor release with non-breaking fixes, enhancements, and new features
* `patch` to perform a patch release with only non-breaking fixes

This script will automatically change the version number (in the `VERSION` file), update the `CHANGELOG.md` file with the changes since the previous release, commit both changes to Git, and tag the last commit (e.g., `v1.0.1` for version 1.0.1).

Review these commit(s) so they are valid, and then push them to the upstream repository:

    $ git push --follow-tags upstream

Next, build the release using this tag:

    $ make all

Verify the functionality works as expected, and then upload the `out/strongback-<version>-<os>.tar.gz` and `out/strongback-<version>-<os>.zip` artifacts to GitHub [as a new release](https://github.com/strongback/strongback-cli/releases). Use the existing tag for the release (e.g., `v1.0.1`) and use the version number as the name the release (e.g., "1.0.1").