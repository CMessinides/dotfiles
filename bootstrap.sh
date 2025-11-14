#!/bin/bash

# Create all directories that should be "real" so stow doesn't make them symlink
# into this repo. Otherwise, files that should stay local could get added to the
# repo and make Git output confusing.
#
# Also sets up some default global Git config because there's one file I want to
# change per machine, so symlinking with stow doesn't make sense.

abbr-home() {
    echo "~${1#$HOME}"
}

_closest-link-in-path() {
    local path=$1
    if [[ "$path" == "$HOME" || "$path" == "." ]]; then
        echo ""
        return
    fi

    local dir=$(dirname "$path")
    local above=$(_closest-link-in-path "$dir")

    if [[ "$above" != "" ]]; then
        echo "$above"
    elif [[ -L "$path" ]]; then
        echo "$path"
    else
        echo ""
    fi
}

_remove-closest-link() {
    local path=$1
    local closest_link=$(_closest-link-in-path "$path")

    if [[ -n "$closest_link" ]]; then
        echo " ! Found symlink at $(abbr-home "$closest_link"); removing..."
        rm "${closest_link%/}"
    fi
}

_ensure-real-dir() {
    local path=$1
    _remove-closest-link "$path"
    if [[ -d "$path" ]]; then
        echo " | Directory already exists at $(abbr-home "$path")"
    else
        mkdir -p "$path"
        echo " > Created directory at $(abbr-home "$path")"
    fi
}

ensure-home-dirs() {
    for path in $@
    do
        _ensure-real-dir "$HOME/$path"
    done
}

ensure-home-dirs .config .local/{bin,share,state}
echo "✅ Base directories"

ensure-home-dirs .config/git
echo "✅ Git directories"

GIT_USER_CONFIG="$HOME/.config/git/user-config"
if ! [[ -f "$GIT_USER_CONFIG" ]]; then
    echo '
    [user]
        name = Cameron Messinides
        email = cameron.messinides@gmail.com

    # vim: ft=gitconfig
    ' > "$HOME/.config/git/user-config"
else
    echo " | Git user config already exists at $(abbr-home "$GIT_USER_CONFIG")"
fi
echo "✅ Git user config"

ensure-home-dirs .config/fish/{completions,conf.d,functions}
ensure-home-dirs .local/share/fish/vendor_completions.d
echo "✅ Fish directories"

ensure-home-dirs .config/nvim/{lua/custom/plugins,luasnippets,queries}
echo "✅ Neovim directories"

ensure-home-dirs .config/tmux/plugins
echo "✅ Tmux directories"

ensure-home-dirs .config/tmuxp
echo "✅ Tmuxp directories"

ensure-home-dirs .doom.d/custom
echo "✅ Doom Emacs directories"
