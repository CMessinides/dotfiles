if command -q bat
    if command -q ddbat
        set -gx DEVDOCS_PAGER "ddbat"
    else
        set -gx DEVDOCS_PAGER "bat"
    end
end
