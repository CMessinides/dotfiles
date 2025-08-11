#!/bin/bash

# Create all directories that should be "real" so stow doesn't make them symlink
# into this repo. Otherwise, files that should stay local could get added to the
# repo and make Git output confusing.

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
        echo " ! Found symlink at $closest_link; removing..."
        rm "${closest_link%/}"
    fi
}

_ensure-real-dir() {
    local path=$1
    _remove-closest-link "$path"
    if [[ -d "$path" ]]; then
        echo " | Directory already exists at $path"
    else
        mkdir -p "$path"
        echo " > Created directory at $path"
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
