# Set fish as the default shell for things like tmux
set -gx SHELL "$(command -s fish)"

# Set editor to the installed flavor of vim
if type -q nvim
    set -gx EDITOR nvim
else
    set -gx EDITOR vim
end

if status is-interactive
    # Commands to run in interactive sessions can go here

    # <Ctrl-F> to open a project launcher
    bind \cf 'open-project; commandline -f repaint'

    # <Ctrl-W> to open a worktree launcher
    bind \cw 'open-worktree; commandline -f repaint'

    # Abbreviations
    abbr -a google search --engine=google
    abbr -a mdn search --engine=mdn
end
