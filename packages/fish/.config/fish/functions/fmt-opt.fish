function fmt-opt -d "Format an option for print-help"
    argparse --max-args 2 \
        'h/help' 's/short=' 'p/placeholder=' -- $argv
    or return

    if set -ql _flag_h
        print-help "Format an option for print-help" \
            --usage "fmt-opt [OPTIONS] <NAME> [DESCRIPTION]" \
            --argument (fmt-arg "NAME" "Long name of the option (without leading --)") \
            --argument (fmt-arg --optional "DESCRIPTION" "Brief help text for the option") \
            --option (fmt-opt "short" \
                        "Short version of the option (without leading -)" \
                        --short "s" \
                        --placeholder "ALIAS") \
            --option (fmt-opt "placeholder" \
                        "Placeholder text for option value" \
                        --short "p" \
                        --placeholder "PLACEHOLDER") \
            --option (fmt-opt "help" "Show this help" --short "h")
        return
    end

    set -l name "$argv[1]"
    set -l description "$argv[2]"

    if test -z "$name"
        print-error "name is required"
        return 1
    end

    if set -ql _flag_s
        if test (string length -- "$_flag_s") -ne 1
            print-error "alias must be a single character"
            return 1
        end

        printf "-%s, " "$_flag_s"
    else
        printf "    "
    end

    printf "--%s" "$name"

    if set -ql _flag_p
        printf " %s" "$_flag_p"
    end

    if test -n "$description"
        printf "\t%s" "$description"
    end

    printf "\n"
end
