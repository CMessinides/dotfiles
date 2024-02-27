export PROJECTS_ROOT="$HOME/source"

open-project-launcher() {
    [ -n "$ZLE_STATE" ] && trap 'zle reset-prompt' EXIT

    local tmuxp_configs=$(tmuxp ls)
    local project=$(find "$PROJECTS_ROOT"/* -maxdepth 0 -type d -exec basename {} \; | cat - <(echo $tmuxp_configs) | sort | uniq | fzf)

    (
        # Restore the standard std* file descriptors for tmux.
        # https://unix.stackexchange.com/a/512979/22339
        exec </dev/tty; exec <&1;
        if [ -n "$project" ]; then
            if echo $tmuxp_configs | grep -x "$project" > /dev/null; then
                tmuxp load "$project"
            else
                local dir="$PROJECTS_ROOT/$project"

                tmux has-session -t "$project" 2>/dev/null
                if [ $? != 0 ]; then
                    tmux new-session -d -c "$dir" -s "$project"
                fi

                if [ -n "$TMUX" ]; then
                    tmux switch-client -t "$project"
                else
                    tmux attach-session -t "$project"
                fi
            fi
        fi
    )
}

zle -N open-project-launcher
bindkey '^F' open-project-launcher
