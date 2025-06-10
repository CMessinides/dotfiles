function open-worktree
    set -f project "$(pick-project)"
    if test -z "$project"
        return 1
    end

    set -f worktree "$(pick-worktree "$project")"
    if test -z "$worktree"
        return 1
    end

    if not test -d "$WORKTREE_ROOT/$project/$worktree"
        return 1
    end

    set -f session_name "$(string replace --all '/' '_' "$project/$worktree")"

    tmux has-session -t "$session_name" &>/dev/null
    or tmux new-session -d -c "$WORKTREE_ROOT/$project/$worktree" -s "$session_name"

    if test -n "$TMUX"
        tmux switch-client -t "$session_name"
    else
        tmux attach-session -t "$session_name"
    end
end
