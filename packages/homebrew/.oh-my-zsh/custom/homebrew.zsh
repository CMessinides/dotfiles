BREW_EXECUTABLES=(
    "/home/linuxbrew/.linuxbrew/bin/brew"
    "/opt/homebrew/bin/brew"
    "/usr/local/bin/brew"
)

for candidate in "${BREW_EXECUTABLES[@]}"
do
    if [ -x "$candidate" ] && [ -f "$candidate" ]; then
        eval "$("$candidate" shellenv)"
        break
    fi
done

unset BREW_EXECUTABLES
