function fmt-arg -d "Format an argument for print-help"
    argparse --max-args 2 \
        'h/help' 'o/optional' 'm/multiple' -- $argv
    or return

    if set -ql _flag_h
        print-help "Format an argument for print-help" \
            --usage "fmt-arg [OPTIONS] <NAME> [DESCRIPTION]" \
            --argument (fmt-arg "NAME" "Name of the argument") \
            --argument (fmt-arg --optional "DESCRIPTION" "Brief help text for the argument") \
            --option (fmt-opt "optional" "Format as optional" --short "o") \
            --option (fmt-opt "multiple" "Format as variadic" --short "m") \
            --option (fmt-opt "help" "Show this help" --short "h")
        return
    end

    set -l name "$argv[1]"
    set -l description "$argv[2]"

    if test -z "$name"
        print-error "name is required"
        return 1
    end

    if set -ql _flag_o
        printf "[%s" $name
    else
        printf "<%s" $name
    end

    if set -ql _flag_m
        printf "..."
    end

    if set -ql _flag_o
        printf "]"
    else
        printf ">"
    end

    if test -n "$description"
        printf "\t%s" "$description"
    end

    printf "\n"
end
