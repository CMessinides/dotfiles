set GO_INSTALL "/usr/local/go"

if [ -d "$GO_INSTALL" ]
    fish_add_path "$GO_INSTALL/bin"
end
