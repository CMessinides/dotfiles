function search
    argparse 'h/help' 'e/engine=?!test -z $_flag_value; or contains $_flag_value $SEARCH_ENGINES' -- $argv
    or return

    if set -ql _flag_help
        echo "search [-h|--help] [-e|--engine=($(string join "|" $SEARCH_ENGINES))] [SEARCH_TERM ...]"
        return 0
    end

    if test -z $_flag_engine
        set -f _flag_engine $SEARCH_ENGINES[1]
    end

    set -l query "$(string escape --style=url $argv | string join '+')"

    switch $_flag_engine
        case google
            open "https://www.google.com/search?q=$query"
        case mdn
            open "https://developer.mozilla.org/en-US/search?q=$query"
    end
end
