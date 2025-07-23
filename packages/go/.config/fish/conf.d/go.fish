set GO_INSTALL "/usr/local/go"
if [ -d "$GO_INSTALL" ]
    fish_add_path "$GO_INSTALL/bin"
end

set -gx GOPATH "$HOME/go"
if command -q go
    fish_add_path "$(go env GOPATH)/bin"
end
