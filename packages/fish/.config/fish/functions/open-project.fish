function open-project -d 'Open a project in tmux'
    set -l project_root
    if set -q PROJECT_ROOT
        set project_root "$PROJECT_ROOT"
    else
        set project_root "$HOME/source"
    end

    set -l project "$(find "$project_root"/* -maxdepth 0 -type d -exec basename {} \; | fzf)"

    if [ -n "$project" ]
        # Test for existing session
        tmux has-session -t "$project" &>/dev/null

        if [ $status != 0 ]
            tmux new-session -d -c "$project_root/$project" -s "$project"
        end

        if [ -n "$TMUX" ]
            tmux switch-client -t "$project"
        else
            tmux attach-session -t "$project"
        end
    end
end
