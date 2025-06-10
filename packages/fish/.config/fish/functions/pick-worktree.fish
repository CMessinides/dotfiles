function pick-worktree
    set -f project "$argv[1]"
    if test -z "$project"
        return 1
    end

    set -f worktree "$(string replace "worktree $WORKTREE_ROOT/$project/" '' (git -C "$PROJECT_ROOT/$project" worktree list --porcelain | grep "^worktree $WORKTREE_ROOT/$project/") | fzf)"
        or return 1

    if test -z "$worktree"
        return 1
    end

    echo "$worktree"
end
