function _search_widget
    set -f engine $(printf '%s\n' $SEARCH_ENGINES | fzf --reverse --prompt '🔎 Search engine: ')
        or return

    set -f terms $(_print_search_history $engine | fzf --scheme history --reverse --prompt "🔎 Search ($engine): " --print-query | tail -1)

    if test -z $terms
        return 1
    end

    _save_search_history $engine $terms
    search --engine=$engine (string split ' ' $terms)
end
