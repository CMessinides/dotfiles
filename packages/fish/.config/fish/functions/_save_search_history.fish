function _save_search_history
    _touch_search_history

    set -l engine $argv[1]
    set -l terms $argv[2..]
    echo ":$engine:" $terms >> $SEARCH_HISTORY_FILE
end
