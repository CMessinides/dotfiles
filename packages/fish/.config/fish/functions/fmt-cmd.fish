function fmt-cmd -d "Format a command for print-help"
    argparse --max-args 2 'h/help' -- $argv
    or return

    if set -ql _flag_h
        print-help "Format a command for print-help" \
            --usage "fmt-cmd [OPTIONS] <NAME> [DESCRIPTION]" \
            --argument (fmt-arg "NAME" "Name of the command") \
            --argument (fmt-arg --optional "DESCRIPTION" "Brief help text for the command") \
            --option (fmt-opt "help" "Show this help" --short "h")
        return
    end

    set -l name "$argv[1]"
    set -l description "$argv[2]"

    if test -z "$name"
        print-error "name is required"
        return 1
    end

    printf "%s" "$name"

    if test -n "$description"
        printf "\t%s" "$description"
    end

    printf "\n"
end
