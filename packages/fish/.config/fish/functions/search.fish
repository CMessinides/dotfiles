set engines google mdn
set default_engine engines[1]

function search
    argparse 'h/help' 'e/engine=?!test -z $_flag_value; or contains $_flag_value $engines' -- $argv
    or return

    if set -ql _flag_help
        echo "search [-h|--help] [-e|--engine=($(string join "|" $engines))] [SEARCH_TERM ...]"
        return 0
    end

    if test -z $_flag_engine
        set -f _flag_engine $default_engine
    end

    set -l escaped_terms ""
    set -l query "$(string escape --style=url $argv | string replace -a ' ' '+')"

    switch $_flag_engine
        case google
            open "https://www.google.com/search?q=$query"
        case mdn
            open "https://developer.mozilla.org/en-US/search?q=$query"
    end
end
