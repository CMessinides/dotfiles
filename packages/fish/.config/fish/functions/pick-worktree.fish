function pick-worktree
    set -f project "$argv[1]"
    if test -z "$project"
        return 1
    end

        set -f multiple 1
    end

    set -f fzf_args
    if test "$argv[2]" = '--multiple'
        set -a fzf_args '-m'
    end

    set -f opts "$(echo '(Create new worktree)'; string replace "worktree $WORKTREE_ROOT/$project/" '' (git -C "$PROJECT_ROOT/$project" worktree list --porcelain | grep "^worktree $WORKTREE_ROOT/$project/"))"

    set -f worktree "$(echo $opts | fzf $fzf_args)"
        or return 1

    if test "$worktree" = '(Create new worktree)'
        set -f worktree "$(cd "$PROJECT_ROOT/$project"; add-worktree)"
            or return 1
    end

    if test -z "$worktree"
        return 1
    end

    echo "$worktree"
end
