#!/bin/bash
#
# Run scripts to setup my developer environment.

#region Bootstrapping

set -o nounset
set -o pipefail
set -o errexit

# Get the absolute directory name from a filename.
# Arguments:
#   $1 - the filename, a path.
# Outputs:
#   Writes absolute directory name to stdout.
absdirname() ( cd "$(dirname "$1")" && pwd )

#endregion

#region Global constants

SRC_ROOT="$(absdirname "${BASH_SOURCE[0]}")"
LIB_DIR="$SRC_ROOT/lib"
COPY_DIR="$SRC_ROOT/copy"

#endregion

#region Helper functions

# Detect whether a commmand exists.
# Arguments:
#   $1 - the command name.
# Outputs:
#   Exits with status 0 if command exists, non-zero status otherwise.
has() {
	command -v "$1" 1>/dev/null 2>&1
}

COLOR_NORMAL=$(tput sgr0)
COLOR_DIM=$(tput dim)
COLOR_GREEN=$(tput setaf 2)
COLOR_YELLOW=$(tput setaf 3)

log() {
	echo "$@" 1>&2
}

log_success() {
	log "${COLOR_GREEN}$@${COLOR_NORMAL}"
}

log_notice() {
	log "${COLOR_YELLOW}$@${COLOR_NORMAL}"
}

log_dim() {
	log "${COLOR_DIM}$@${COLOR_NORMAL}"
}

install() {
    if has brew; then
        log_dim "Homebrew is already installed"
    else
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        log_success "Homebrew installed"
    fi

    brew bundle install --no-lock --file "$SRC_ROOT/brew/.Brewfile"
    log_success "Homebrew packages installed"

    local zsh_path="$(command -v zsh)"
    if [ "$zsh_path" = "$SHELL" ]; then
        log_dim "zsh is already the default shell"
    else
        chsh -s "$zsh_path"
        log_success "zsh set as default shell"
    fi

    if [ -d "$HOME/.oh-my-zsh/" ]; then
        log_dim "oh-my-zsh is already installed"
    else
        git clone https://github.com/ohmyzsh/ohmyzsh.git "$HOME/.oh-my-zsh"
        log_success "oh-my-zsh installed"
    fi

    local stow_packages="$(basename -a $COPY_DIR/*)"
    stow -t "$HOME" -d "$COPY_DIR" -S $stow_packages
    log_success "Dotfiles linked"
}

$@
