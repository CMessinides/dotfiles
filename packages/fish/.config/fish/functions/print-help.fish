function print-help -d "Print formatted CLI help text to stderr"
    argparse --max-args 1 \
        'h/help' 'u/usage=' 'a/argument=+' 'o/option=+' -- $argv
    or return

    if set -ql _flag_h
        print-help "Print formatted CLI help text to stderr" \
            --usage "print-help <SUMMARY> [OPTIONS]" \
            --argument (fmt-arg "SUMMARY" "Brief description of the command") \
            --option (fmt-opt "usage" \
                        "Provide one-line command usage" \
                        --short "u" \
                        --placeholder "USAGE") \
            --option (fmt-opt "argument" \
                        "Document an argument. Can be provided for multiple arguments." \
                        --short "a" \
                        --placeholder "ARG") \
            --option (fmt-opt "option" \
                        "Document an option. Can be provided for multiple options." \
                        --short "o" \
                        --placeholder "OPT") \
            --option (fmt-opt "help" "Show this help" --short "h")
        return
    end

    set -l summary "$argv[1]"
    if test -z "$summary"
        print-error "summary is required"
        return 1
    end

    begin
        printf "%s\n" "$summary"

        if set -ql _flag_u
            printf "\nUsage: %s\n" "$_flag_u"
        end

        set -l nargs (count $_flag_a)
        if test $nargs -gt 0
            printf "\nArguments:\n"
            for arg in $_flag_a
                echo "  $arg"
            end | column -ts \t
        end

        set -l nopts (count $_flag_o)
        if test $nopts -gt 0
            printf "\nOptions:\n"
            for opt in $_flag_o
                echo "  $opt"
            end | column -ts \t
        end
    end >&2
end
