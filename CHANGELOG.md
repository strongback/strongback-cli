## 1.2.1
* Bump version to 1.2.1
* Corrected behavior of `--overwrite` flag of `new-project` (#10)

## 1.2.0
* Bump version to 1.2.0
* Expose more information about 3rd party dependencies
* Added Travis-CI automated builds
* Fixed the `new-project` usage to have correct/consistent option names (Issue #5)
* Added badges to README
* Changed the contributing description so that issue numbers are at the end of the first line for each commit message.
* Merge pull request #4 from rhauch/issue-1
* Issue 1 - Make sure archive directory exists
* Merge pull request #3 from rhauch/issue-2
* Changed the Makefile to run `make all` by default
* Issue 2 - Removed developer’s GOPATH from stack trace output

## 1.1.0
* Bump version to 1.1.0
* Added support for copying JARs into WPILib's user library directory
* Changed to building a ZIP file for Windows, since Windows won't extract a .tar.gz file
* Renamed file containing instructions for releasing
* Minor corrections and cleanup

## 1.0.1
* Bump version to 1.0.1
* Corrected output of team number
* Changed build to no longer use ZIP files
* Corrected release scripts

## 1.0.0
* Updated release scripts
* Updated new-project to use different default package and to fail if not specified and WPILib is not yet initialized
* Added printing of team number under WPILib section
* Additional fixes and improvements. Tested on MacOS and Windows.
* Separated the Linux and OS X sections in the README
* Updated the README file
* Initial codebase with complete functionality
