function add-worktree
    if not set -f git_dir "$(git rev-parse --git-dir)"
        return 1
    end

    if test "$git_dir" = ".git"
        set -f project "$(pwd | xargs basename)"
    else
        set -f project "$(dirname $repo_root | xargs basename)"
    end

    set -l branch "$argv[1]"

    if test -z "$branch"
        if not set branch "$(string trim (git branch -a | grep -v '^[*+]') | fzf)"
            return 1
        end
    end

    set -l worktree_path "$WORKTREE_ROOT/$project/$branch"
    if not git worktree add "$worktree_path" "$branch" >&2
        return 1
    end

    echo "$worktree_path"
end
