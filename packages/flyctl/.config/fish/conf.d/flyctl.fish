set FLYCTL_INSTALL "$HOME/.fly"

if [ -d "$FLYCTL_INSTALL" ]
    fish_add_path "$FLYCTL_INSTALL/bin"
end
