function open-project -d 'Open a project in tmux'
    set -l project "$(pick-project)"

    if [ -n "$project" ]
        # Test for existing session
        tmux has-session -t "$project" &>/dev/null

        if [ $status != 0 ]
            tmux new-session -d -c "$PROJECT_ROOT/$project" -s "$project"
        end

        if [ -n "$TMUX" ]
            tmux switch-client -t "$project"
        else
            tmux attach-session -t "$project"
        end
    end
end
