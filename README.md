# Dotfiles

These are my configuration files.

## Installation

Clone this repository. By default, the installation script (see below) expects
to find the repo at `~/dotfiles`.

```shell
$ git clone https://github.com/cmessinides/dotfiles ~/dotfiles
```

Run the [installation script](./install.sh) to install dependencies and link
the dotfiles into the correct locations.

```shell
$ ~/dotfiles/install.sh
```

If you just want to see what the script would do, without actually making any
changes, add the `--dry-run` flag:

```shell
$ ~/dotfiles/install.sh --dry-run
```

If you cloned the repository to a different location (not `~/dotfiles`), you'll
need to set the `DOTFILES` variable to the correct directory:

```shell
$ DOTFILES="~/.config/dotfiles" ~/.config/dotfiles/install.sh
```
