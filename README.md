# releaseman

Your friendly Release helper.

Using this tool is as easy as running `releaseman create` and following the guide
it prints. `releaseman` helps you generating changelog and tagging your states.

**What this tool does:**

1. Collects your git commit messages between start and end states
2. Creates changelog based on commit messages
3. Adds the changelog to your git and makes commit on development branch
4. Merges your development branch into your release branch and commits
5. Tag the release branch with the release version

**What this tool does not:**

1. **Does not push into your git repository, so you can roll back all the changes**

*Roll back:*

* if you want to undo the last commit you can call:
  `git reset --hard HEAD~1`
* to roll back to the remote state:'
  `git reset --hard origin/[branch-name]`

### How to use

1. Simply you can start creating release with `releaseman create` and following the guide
2. You can provide required inputs as cli arguments `releaseman create --development-branch develop --release-branch master --start-state 1.0.0 --end-state "current state" --changelog-path ./changelog --release-version 1.0.1`
3. Or you can provide a configuration file in your repository root (where you run releaseman) with name `config.yml` format example: THIS_REPO_ROOT/configs/config.yml
4. If you provide config.yml and also cli params, cli params will override the values in the config.yml (only during the run). This allows you to provide static params in your config, and dynamic ones as cli args.
