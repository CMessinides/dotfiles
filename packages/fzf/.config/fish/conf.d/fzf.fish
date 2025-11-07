# Set up fzf key bindings
fzf --fish | source

# Prettify the default fzf
set -gx FZF_DEFAULT_OPTS "
    --border horizontal --reverse --cycle
    --preview-window 'right,border-left,<66(down,50%,border-top)'
    --color 16
    --color border:grey
    --color prompt:blue
    --color pointer:blue
    --color info:italic:blue
    --color label:italic:blue
    --color fg+:bold:bright-magenta,bg+:-1
    --color hl:underline:yellow,hl+:bold:underline:bright-yellow
    --prompt '❱ '
    --pointer '▒'
"
