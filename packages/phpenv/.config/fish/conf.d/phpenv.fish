set -gx PHPENV_ROOT "$HOME/.phpenv"

if [ -d "$PHPENV_ROOT" ]
  fish_add_path "$PHPENV_ROOT/bin"
  eval "$(phpenv init -)"
end
