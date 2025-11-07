function pick-project -d "Pick a project with fzf"
    find "$PROJECT_ROOT"/* -maxdepth 0 -type d -exec basename {} \; | fzf \
        --border-label " Projects " \
        --border-label-pos 2 \
        --height '~25%'
end
