#!/usr/bin/env bash
# :help:
# Setup my developer environment on a new machine.
#
# Usage:
# 	~/dotfiles/install.sh [--dry-run] [--help]
#
# Options:
# 	--dry-run		Don't make any changes; just explain them.
# 	-h, --help		Display this help.
# 
# Customization:
#
# This script expects that the dotfiles repository has been cloned to
# $HOME/dotfiles. If the repo was cloned elsewhere, let the script know by
# setting the DOTFILES variable.
#
# 	DOTFILES="~/my-dotfiles" ~/my-dotfiles/install.sh
# :endhelp:
set -o nounset
set -o pipefail
set -o errexit

eval_unsafe() {
	set +o nounset
	set +o pipefail
	set +o errexit
	eval $@
	set -o nounset
	set -o pipefail
	set -o errexit
}

# Customizable variables
# DOTFILES: Location of the cloned dotfiles repository
DOTFILES="${DOTFILES:-$HOME/dotfiles}"

# Internal state
DRY_RUN=0
NUM_SUGGESTED_CHANGES=0

is_dry_run() {
	[ $DRY_RUN -ne 0 ]
}

has_suggested_changes() {
	[ $NUM_SUGGESTED_CHANGES -gt 0 ]
}

add_suggested_change() {
	NUM_SUGGESTED_CHANGES=$((NUM_SUGGESTED_CHANGES+1))
}

normal=$(tput sgr0)
dim=$(tput dim)
green=$(tput setaf 2)
yellow=$(tput setaf 3)

log() {
	echo "$@" 1>&2 
}

log_success() {
	log "${green}$@${normal}"
}

log_notice() {
	log "${yellow}$@${normal}"
}

log_dim() {
	log "${dim}$@${normal}"
}

has() {
	command -v "$1" 1>/dev/null 2>&1
}

help() {
	sed -n '/^# :help:/,/^# :endhelp:/{
			/^# :help:/!{
				/^# :endhelp:/!p
			};
			/^# :endhelp:/q
		}' $BASH_SOURCE | cut -c 3-
}

_get_zsh_path() {
	command -v zsh
}

_safe_install_zsh() {
	if is_dry_run
	then
		log_notice "Dry run: would have installed zsh"
		log_dim '└ $ sudo apt install zsh'
		add_suggested_change
	else
		log "Installing zsh..."
		sudo apt install zsh
		log_success "Installed zsh"
	fi
}

_safe_set_zsh_as_default() {
	local zsh_path=$1
	if is_dry_run
	then
		log_notice "Dry run: would have set zsh as the default shell"
		log_dim "└ $ chsh -s \"${zsh_path:-<unknown>}\""
		add_suggested_change
	else
		log "Setting zsh as the default shell..."
		chsh -s "$zsh_path"
		log_success "Default shell is zsh"
	fi
}

_safe_install_oh_my_zsh() {
	if is_dry_run
	then
		log_notice "Dry run: would have installed oh-my-zsh"
		log_dim '└ $ git clone https://github.com/ohmyzsh/ohmyzsh.git ~/.oh-my-zsh'
		add_suggested_change
	else
		log "Installing oh-my-zsh..."
		git clone https://github.com/ohmyzsh/ohmyzsh.git $HOME/.oh-my-zsh
		log_success "Installed oh-my-zsh"
	fi
}

_safe_install_starship() {
	if is_dry_run
	then
		log_notice "Dry run: would have installed Starship"
		log_dim '└ $ curl -sS https://starship.rs/install.sh | sh'
		add_suggested_change
	else
		log "Installing starship..."
		curl -sS https://starship.rs/install.sh | sh
		log_success "Installed starship"
	fi
}

_safe_install_fnm() {
	if is_dry_run
	then
		log_notice "Dry run: would have installed fnm"
		log_dim '└ $ curl -fsSL https://fnm.vercel.app/install | bash -s -- --skip-shell'
		add_suggested_change
	else
		log "Installing fnm..."
		curl -fsSL https://fnm.vercel.app/install | bash -s -- --skip-shell
		log_success "Installed fnm"
	fi
}

_safe_install_node() {
	if is_dry_run
	then
		log_notice "Would have installed node"
		log_dim "└ $ fnm install --lts"
		add_suggested_change
	else
		log "Installing node..."
		fnm install --lts
		log_success "Installed node"
	fi
}

_safe_install_deno() {
	if is_dry_run
	then
		log_notice "Dry run: would have installed Deno"
		log_dim '└ $ curl -fsSL https://deno.land/x/install/install.sh | sh'
		add_suggested_change
	else
		log "Installing Deno..."
		curl -fsSL https://deno.land/x/install/install.sh | sh
		log_success "Installed Deno"
	fi
}

_safe_install_rust() {
	if is_dry_run
	then
		log_notice "Dry run: would have installed Rust"
		log_dim '└ $ curl --proto "=https" --tlsv1.2 -sSf https://sh.rustup.rs | sh'
		add_suggested_change
	else
		log "Installing Rust..."
		curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
		log_success "Installed Rust"
	fi
}

_safe_install_neovim() {
	if is_dry_run
	then
		log_notice "Dry run: would have installed neovim"
		log_dim '├ $ sudo apt-get install software-properties-common'
		log_dim '├ $ sudo add-apt-repository ppa:neovim-ppa/stable'
		log_dim '├ $ sudo apt-get update'
		log_dim '└ $ sudo apt-get install neovim'
		add_suggested_change
	else
		log "Installing neovim..."
		sudo apt-get install software-properties-common
		sudo add-apt-repository ppa:neovim-ppa/stable
		sudo apt-get update
		sudo apt-get install neovim
		log_success "Installed neovim"
	fi
}

_safe_install_system_package() {
	local package="$1"
	if is_dry_run
	then
		log_notice "Dry run: would have installed $package"
		log_dim "└ $ sudo apt-get install $package"
		add_suggested_change
	else
		log "Installing $package..."
		sudo apt-get install $package
		log_success "Installed $package"
	fi
}

_safe_install_npm_package() {
	local package="$1"
	if is_dry_run
	then
		log_notice "Dry run: would have installed $package"
		log_dim "└ $ npm install --global \"$package\""
		add_suggested_change
	else
		log "Installing $package..."
		npm install --global "$package"
		log_success "Installed $package"
	fi
}

_safe_stow_dotfiles() {
	local packages="git nvim fnm starship tmux tmuxp zsh"
	if is_dry_run
	then
		log_notice "Dry run: would have linked dotfiles with stow"
		log_dim "└ $ stow -t \"$HOME\" -d \"$DOTFILES\" -S $packages"
		log $dim
		stow --simulate --verbose -t "$HOME" -d "$DOTFILES" -S $packages
		log $normal
		add_suggested_change
	else
		log "Stowing dotfiles..."
		stow -t "$HOME" -d "$DOTFILES" -S $packages
		log_success "Stowed dotfiles"
	fi
}

_install_system_package_if_needed() {
	local package="$1"
	local bins="${2:-$package}"
	local is_missing_bins=0

	for bin in $bins
	do
		if ! has "$bin"
		then
			is_missing_bins=1
			break
		fi
	done

	if [ $is_missing_bins -eq 0 ]
	then
		log_dim "$package is already installed"
	else
		_safe_install_system_package "$package"
	fi
}

_install_npm_package_if_needed() {
	local package="$1"
	local bins="${2:-$package}"
	local is_missing_bins=0

	for bin in $bins
	do
		if ! has "$bin"
		then
			is_missing_bins=1
			break
		fi
	done

	if [ $is_missing_bins -eq 0 ]
	then
		log_dim "$package is already installed"
	else
		_safe_install_npm_package "$package"
	fi
}

install_shell() {
	# Install zsh
	if has zsh
	then
		log_dim "zsh is already installed"
	else
		_safe_install_zsh
	fi

	# Set zsh as default shell
	local zsh_path="$(_get_zsh_path)"
	if [ "$zsh_path" = "$SHELL" ]
	then
		log_dim "zsh is already the default shell"
	else
		_safe_set_zsh_as_default "$zsh_path"
	fi

	# Install oh-my-zsh
	if [ -d "$HOME/.oh-my-zsh/" ]
	then
		log_dim "oh-my-zsh is already installed"
	else
		_safe_install_oh_my_zsh
	fi

	# Install Starship prompt
	if has starship
	then
		log_dim "starship is already installed"
	else 
		_safe_install_starship
	fi
}

install_dotfiles() {
	# Install stow (for managing dotfiles)
	_install_system_package_if_needed stow
	
	# Sync the dotfiles with stow
	_safe_stow_dotfiles
}

install_developer_tools() {
	# Install Fast Node Manager (fnm)
	if has fnm
	then
		log_dim "fnm is already installed"
	else
		_safe_install_fnm
	fi

	if [ "$(fnm current)" != "none" ]
	then
		log_dim "node is already installed"
	else
		_safe_install_node
	fi

	# Install Deno
	if has deno
	then
		log_dim "Deno is already installed"
	else
		_safe_install_deno
	fi

	# Install Rust
	if has rustc
	then
		log_dim "Rust is already installed"
	else
		_safe_install_rust
	fi

	# Install neovim
	if has nvim
	then
		log_dim "neovim is already installed"
	else
		_safe_install_neovim
	fi

	# Install ripgrep (neovim-telescope dependency)
	_install_system_package_if_needed ripgrep rg

	# Install build-essential (various neovim dependencies)
	_install_system_package_if_needed build-essential "make gcc"

	# Install tmuxp (tmux workspace manager)
	_install_system_package_if_needed tmuxp
}

install_utilities() {
	# Install jq
	_install_system_package_if_needed jq

	# Install unzip
	_install_system_package_if_needed unzip
}

install() {
	install_shell
	install_dotfiles
	install_utilities
	install_developer_tools

	if has_suggested_changes
	then
		log_notice "To run the above steps, remove the --dry-run flag and rerun this script."
	else
		log_success "Setup complete. Enjoy!"
	fi
}

for arg in "$@"
do
	case "$arg" in
	-h)			help
				exit
				;;
	--help)		help
				exit
				;;
	--dry-run)	DRY_RUN=1
				;;
	*)			log "Unrecognized option: $arg"
				exit 1
				;;
	esac
done

install
