# gh-automerge

Github PRs auto merger tool.

## How to install

    go get -u github.com/nordicdyno/gh-automerge

## How to setup

### Step1. Get Github Oauth Token

Create it here: https://github.com/settings/tokens (only `repo/public_repo` checkbox required) and save it somewhere.

### Step2. Setup local environment

put to your `.bashrc` file string:

    export GITHUB_AUTH_TOKEN=<token-value>
    export GITHUB_PROJECT=<default-repository-owner>
    export GITHUB_REPO=<default-repository-name>

or pass token via `-token`, `-p` and `-r` flags to gh-automerge command.

## How to use

with properly defined environment variables, just pass pull request number:

    gh-automerge -pr=<PR-NUMBER>
