set -l commands capture edit list ls remove rm show

function __fish_notes_needs_command
    set -l cmd (commandline -xpc)

    if test (count $cmd) -eq 1
        return 0
    end

    return 1
end

function __fish_notes_using_command
    set -l cmd (commandline -xpc)

    if test (count $cmd) -gt 1
        if test $argv[1] = $cmd[2]
            return 0
        end
    end

    return 1
end

complete -c notes -f

# commands
complete -c notes -n __fish_notes_needs_command -a capture -d "Capture a quick note"
complete -c notes -n __fish_notes_needs_command -a edit -d "Edit a new or existing note"
complete -c notes -n __fish_notes_needs_command -a "list ls" -d "List all notes"
complete -c notes -n __fish_notes_needs_command -a "remove rm" -d "Delete notes"
complete -c notes -n __fish_notes_needs_command -a "show" -d "Print a note's content"

# capture
# target options
complete -x -c notes -n '__fish_notes_using_command capture' -o t -l title -a "(notes list)" -d "Set title for the captured note"
complete -c notes -n '__fish_notes_using_command capture' -o d -l daily -d "Save capture to the daily note"
complete -c notes -n '__fish_notes_using_command capture' -o w -l weekly -d "Save capture to the weekly note"
complete -c notes -n '__fish_notes_using_command capture' -o m -l monthly -d "Save capture to the monthly note"
complete -c notes -n '__fish_notes_using_command capture' -o y -l yearly -d "Save capture to the yearly note"
# source options
complete -r -c notes -n '__fish_notes_using_command capture' -o f -l from-file -F -d "Init capture with content from file"
