# Set up fzf key bindings
fzf --fish | source

# Prettify the default fzf
set -gx FZF_DEFAULT_OPTS "
    --border horizontal --reverse --cycle
    --preview-window 'right,border-left,<66(down,50%,border-top)'
    --color 16
    --color border:grey
    --color prompt:bright-magenta
    --color pointer:bright-blue
    --color gutter:grey
    --color info:italic:magenta
    --color label:italic:magenta
    --color fg+:bold:bright-blue,bg+:grey
    --color hl:underline:green,hl+:bold:underline:bright-green
    --prompt '❱ '
"
