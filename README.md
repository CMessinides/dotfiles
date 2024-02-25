# Dotfiles

These are my configuration files. Not really meant to be forked, as these are
all my personal preference, but feel free to copy and take inspiration.

## About

Dotfiles are arranged in packages that can be installed and uninstalled
individually using a symlink manager like stow (see below).

## Requirements

- Python 3.10+
- [GNU Stow](https://www.gnu.org/software/stow/) 2.3.1

## Installation

Clone this repository to a location in your filesystem (`~/dotfiles` is
recommended).

```sh
$ git clone https://github.com/cmessinides/dotfiles ~/dotfiles
```
[`home-manager.py`](./home-manager.py) is a little CLI written in Python that
provides some helpful commands for managing the dotfiles.

List the available packages:

```sh
$ ./home-manager.py list
```

Get the status of a package:

```sh
$ ./home-manager.py status zsh
```

Install one or more packages by name:

```sh
$ ./home-manager.py install zsh nvim
```

Install all packages:

```sh
$ ./home-manager.py install --all
```

Uninstall a package:

```sh
$ ./home-manager.py uninstall flyctl
```

Uninstall all packages:

```sh
$ ./home-manager.py uninstall --all
```
