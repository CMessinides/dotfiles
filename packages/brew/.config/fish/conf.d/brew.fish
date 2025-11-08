set -l BREW_HOME "/home/linuxbrew/.linuxbrew"
set -l BREW_COMMAND "$BREW_HOME/bin/brew"

if test -x "$BREW_COMMAND"
    eval "$("$BREW_COMMAND" shellenv)"
end
