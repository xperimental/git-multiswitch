# git-multiswitch

Command-line tool to switch multiple git repositories contained in a directory to the same branch.

## Install

If you have a recent (>= 1.19) Go installation and your `$GOPATH/bin` directory is in your `$PATH` then this should be enough to get the tool installed:

```shell
go install github.com/xperimental/git-multiswitch@latest
```

If you want to checkout the source and customize your build:

```shell
# Clone the repository and switch to directory
git clone https://github.com/xperimental/git-multiswitch.git
cd git-multiswitch
# Build binary
make
```

Once the `git-multiswitch` binary is in a directory which is in the `$PATH` it is accessible using `git multiswitch` similar to other git commands.

## Usage

```plain
Usage of ./git-multiswitch:
      --base-path string     Contains the path used as starting point for finding projects.
  -b, --branch string        Name of branch to switch to.
      --git-command string   git executable to use.
```

The idea of this tool is to run it in a directory containing git repositories which should be switched to a common branch:

```shell
git multiswitch -b feature-b
```

This will switch all git repositories contained in the current directory to the "feature-b" branch, if a local branch with that name already exists.

### When run from inside a git repository

If you run `git-multiswitch` from inside a repository it will by default "search down" to find repositories to switch to a common branch. This behaviour can be changed using the `--escape-repo` option. When this option is specified and the tool is run from inside a repository, it will first look for the parent directory of the git repository and start its search from there:

```shell
git multiswitch --escape-repo -b feature-b
```
