function add-worktree
    if not set -f git_dir "$(git rev-parse --git-dir)"
        return 1
    end

    if test "$git_dir" = ".git"
        set -f project "$(basename (pwd))"
    else
        set -f project "$(basename (dirname "$git_dir"))"
    end

    set -f branch "$argv[1]"

    if test -z "$branch"
        set -l opts "$(echo '(Create new branch)'; string trim (git branch -a | grep -v '^[*+]'))"
        set -f branch "$(echo $opts | fzf)"

        if test "$branch" = '(Create new branch)'
            read -f branch --prompt-str='Branch name: '
                or return 1
            if test -n "$branch"
                git branch -c "$branch"
            end
        end
    end

    if test -z "$branch"
        return 1
    end

    set -l worktree_path "$WORKTREE_ROOT/$project/$branch"
    if not git worktree add "$worktree_path" "$branch" >&2
        return 1
    end

    echo "$branch"
end
