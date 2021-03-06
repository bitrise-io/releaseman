## Changelog (Current version: 0.9.4)

-----------------

### 0.9.4 (2017 Jan 19)

* [cb7b754] version fix

### 0.9.3 (2017 Jan 19)

* [30ae2fc] debug log
* [40630bd] "development_branch" changed to "master"
* [d801bb3] README: better "quick start guide" section
* [61319f6] go deps update + bitrise.yml for better auto upd (#12)
* [f1c2052] gitignore .gows.user.yml
* [a69e883] Merge pull request #4 from bitrise-tools/godep-update
* [5615f55] updated CI/tests config for Go 1.6+
* [41bbbeb] migrate to vendor/

### 0.9.2 (2016 Mar 02)

* [d718a67] FIX: git tagged commit sort order

### 0.9.1 (2016 Mar 02)

* [4fb0369] bitrise.yml updates
* [a945277] release fix
* [50f8323] Merge develop into master, release: v0.9.1
* [b2d1553] v0.9.1
* [2f74287] logging fix
* [3f198a1] create release on develop
* [f2ef38d] Merge branch 'develop'
* [10b857d] collecting git commits fix
* [1eb6660] commit range fix
* [e9d8904] commit range fix
* [fd19c23] release_config fix
* [abfbb53] test fixes
* [753c5ce] removed changelog
* [8636a1a] release_config fix
* [0fb3de7] use release_config's header and footer templates, commit range fix
* [94a36e7] default relese branch: master

### 0.9.0 (2016 Feb 19)

* [67600e0] Merge develop into master, release: v0.9.0
* [b0d011c] v0.9.0
* [593e250] default values at user inputs
* [3f37a5a] test fixes
* [43c9f60] create custom release version workflow added
* [1ddc416] godeps update
* [fd2ecac] bump last version's patch part by default
* [8e19279] NEW: changelog header and footer support
* [a7c5ee6] FIX: parsing git commits
* [43b951d] FIX: omit not sem-ver git tags
* [972874b] FIX: instead of function name trimm  -> firstChars

### 0.0.4 (2016 Feb 17)

* [c70a0de] Merge develop into master, release: v0.0.4
* [8e784af] v0.0.4
* [df4e2db] set current version to 0.0.3
* [c3abd06] temp rollback
* [a8d3f23] FIX: get changelist
* [839b826] tmp
* [0dbfd98] set -e in release workflow
* [f2c5a74] v0.0.4
* [a5770b2] FIX: add all changes during release
* [2babab9] gitignore _bin/
* [9b9e596] create-release workflow in full release mode
* [8d4428a] create-release now creates binaries too
* [ffb0a10] set current version 0.0.3
* [8018562] NEW: get version / set version script support
* [ce74a44] releaseconfig changes
* [964564d] New input flag: bump-version-script
* [5c5dbcd] refactores
* [be724a9] CI wired in
* [f8af09c] use yml parser instead of templates for config handling
* [6b78b71] Trimm function in changelog template
* [9774d23] ask for next version message
* [085e750] godep update
* [829a841] bitrise ci wired in
* [f782255] ver 0.0.3

### 0.0.3 (2016 Feb 16)

* [b738dee] Merge branch 'master' of github.com:bitrise-tools/releaseman
* [35d931f] Merge pull request #3 from godrei/master
* [7f75607] Merge branch 'develop'
* [f58922e] release fixes
* [96d5a94] v0.0.3
* [e188d8e] Merge develop into master, release: v0.0.3
* [92f23ad] default template contains more options
* [8fa62ba] Added default changelog template to release_config
* [d90db93] Remove whitespace from changelog list
* [e63942f] Remove section separator from default changelog_template
* [5806fdc] Reverse collected commits order
* [69a251f] Allow create changelog with dirty git
* [45f9574] FIX: Only print previous release version
* [9ac525c] revisions: CHANGELOG.md moved; version separated into a version package; fixed pathutil package links

### 0.0.2 (2016 Feb 09)

* [e97a9d7] Merge pull request #2 from godrei/master
* [a306ead] Merge develop into master, release: v0.0.2
* [c1b132e] v0.0.2
* [85faaea] FIX: release config
* [3ed6cd1] init release configurations
* [acefed0] cleanup
* [250cf26] FIX: test fixes
* [7ea5fc3] prepare for new version
* [695967b] code review, some refractors
* [765b440] godep save
* [c5c9687] go test added
* [32b0c3b] FIX: init save template path too
* [b9d13d9] new changelog template
* [539ce94] NEW: separated commands for creating changelog and release
* [f9a18e2] author and commit date added to changelog
* [8e4a13b] list already exist tags befor asking for next version, fail if next version already exist
* [da584a7] NEW: changelog template support
* [0795d1c] release config changes, changelog moved to _changelog
* [d86daf6] create release config and changelog with go template
* [1c3191d] FIX: changelog generation
* [ba58d36] go tests added

-----------------

Updated: 2017 Jan 19