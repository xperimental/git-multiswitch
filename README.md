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

Will switch all repositories contained in the current path to the "feature-b" branch, if a local branch with that name already exists.
