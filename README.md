# releaseman

Your friendly Release helper.

Using this tool is as easy as running `releaseman create` and following the guide
it prints. `releaseman` helps you generating changelog and releasing new version.

**What this tool does:**

1. Generates changelog `releaseman create-changelog`
2. Release new version `releaseman create-release`
3. Generates changelog and release new version `releaseman create`

**What this tool does not:**

1. **Does not push into your git repository, so you can roll back all the changes**

*Roll back:*

* if you want to undo the last commit you can call:
  `git reset --hard HEAD~1`
* to delete tag:
  `$ git tag -d [TAG]`
  `$ git push origin :refs/tags/[TAG]`
* to roll back to the remote state:'
  `git reset --hard origin/[branch-name]`

## How to use

### Init

*Interactive:*

Just start with creating your release configuration, type in `releaseman init` and follow the printed guide.

### Create changelog

*Interactive:*

Type in `releaseman create-changelog` and follow the printed guide.

---

*cli:*

Releaseman needs the following informations to creating changelog:

* `--development-branch`: changelog will generated based on this branchs commits
* `--version`: your current state will marked with this version
* `--bump-version`: if you have tagged git states, use this to auto increment latest tag, and use to mark the current state in changelog [options: patch, minor, major]
* `--changelog-path`

*Evrey input you provide with flag will used instead of the value you provided in your release_config.yml. If you want to use value from config just omitt the related flag.*

May your command looks like:

* `releaseman create-changelog --version 1.1.1` *in common case: use config, and define the missing inut*
* `releaseman create-changelog --bump-version major` *in common case: use config, and define the missing inut, if you have tags*
* `releaseman create-changelog --development-branch develop --version 1.1.1 --changelog-path ./changelog.md` *to override all your configs*
* `releaseman create-changelog --development-branch develop --bump-version major --changelog-path ./changelog.md` *to override all your configs, if you have tags*

---


### Release new version

*Interactive:*

Type in `releaseman create-release` and follow the printed guide.

---

*cli:*

Releaseman needs the following informations for releasing new version:

* `--development-branch`: changes on this branch will merged to release branch
* `--release-branch`: changes on development branch will merged into this branch and this branch will tagged with the release version
* `--version`: release version
* `--bump-version`: if you have tagged git states, use this to auto increment latest tag, and use as release version

*Evrey input you provide with flag will used instead of the value you provided in your release_config.yml. If you want to use value from config just omitt the related flag.*

May your command looks like:

* `releaseman create-release --version 1.1.1` *in common case: use config, and define the missing inut*
* `releaseman create-release --bump-version major` *in common case: use config, and define the missing inut, if you have tags*
* `releaseman create-release --development-branch develop --release-branch master --version 1.1.1` *to override all your configs*
* `releaseman create-release --development-branch develop --release-branch master --bump-version major` *to override all your configs, if you have tags*

---

### Create changelog and Release new version

*Interactive:*

Type in `releaseman create` and follow the printed guide.

---

*cli:*

Releaseman needs the following informations for create changelog and release new version:

* `--development-branch`: changelog will generated based on this branchs commits and changes on this branch will merged to release branch
* `--release-branch`: changes on development branch will merged into this branch and this branch will tagged with the release version
* `--version`: release version
* `--bump-version`: if you have tagged git states, use this to auto increment latest tag, and use as release version
* `--changelog-path`

*Evrey input you provide with flag will used instead of the value you provided in your release_config.yml. If you want to use value from config just omitt the related flag.*

May your command looks like:

* `releaseman create --version 1.1.1` *in common case: use config, and define the missing inut*
* `releaseman create --bump-version major` *in common case: use config, and define the missing inut, if you have tags*
* `releaseman create --development-branch develop --release-branch master --version 1.1.1 --changelog-path ./changelog.md` *to override all your configs*
* `releaseman create --development-branch develop --release-branch master --bump-version major --changelog-path ./changelog.md` *to override all your configs, if you have tags*

---
