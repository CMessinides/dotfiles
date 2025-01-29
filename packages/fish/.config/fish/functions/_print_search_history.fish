function _print_search_history
    _touch_search_history

    set -l engine $argv[1]
    cat $SEARCH_HISTORY_FILE | grep -e "^:$engine: " | cut -d ' ' -f2- | sort | uniq
end
