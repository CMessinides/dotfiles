if command -q bat
    if command -q ddbat
        set -gx DEVDOCS_PAGER "ddbat"
    else
        set -gx DEVDOCS_PAGER "bat"
    end
end

if status is-interactive
    # Search DevDocs documentation
    bind ctrl-s,ctrl-s 'docsearch; commandline -f repaint'

    # Shortcuts for frequently referenced docs
    # Go
    bind ctrl-s,ctrl-g 'docsearch go; commandline -f repaint'
    # CSS
    bind ctrl-s,ctrl-c 'docsearch css; commandline -f repaint'
    # HTML
    bind ctrl-s,ctrl-h 'docsearch html; commandline -f repaint'
    # JavaScript
    bind ctrl-s,ctrl-j 'docsearch javascript; commandline -f repaint'
    # Node.js
    bind ctrl-s,ctrl-n 'docsearch node; commandline -f repaint'
end
