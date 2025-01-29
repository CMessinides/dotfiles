function _touch_search_history
    set -l search_history_dir $(dirname $SEARCH_HISTORY_FILE)

    if not test -d $search_history_dir
        mkdir -p $search_history_dir
    end

    if not test -f $SEARCH_HISTORY_FILE
        touch $SEARCH_HISTORY_FILE
    end
end
