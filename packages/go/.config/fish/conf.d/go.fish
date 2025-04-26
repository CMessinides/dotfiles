set GO_INSTALL "/usr/local/go"
set GOPATH "$HOME/go"

if [ -d "$GO_INSTALL" ]
    fish_add_path "$GO_INSTALL/bin"
    fish_add_path "$GOPATH/bin"
    fish_add_path "$(go env GOPATH)/bin"
end
