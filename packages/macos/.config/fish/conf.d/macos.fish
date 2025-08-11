if command -q brew
    set -gx HOMEBREW_PREFIX "$(brew --prefix)"

    fish_add_path "$HOMEBREW_PREFIX/bin"
    fish_add_path "$HOMEBREW_PREFIX/opt/findutils/libexec/gnubin"
    fish_add_path "$HOMEBREW_PREFIX/opt/gnu-sed/libexec/gnubin"
end
