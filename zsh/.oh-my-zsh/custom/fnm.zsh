export FNM_INSTALL="$HOME/.local/share/fnm"
export PATH="$FNM_INSTALL:$PATH"
eval "$(fnm env --use-on-cd)"
