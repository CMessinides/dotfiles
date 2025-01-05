# fnm
set FNM_PATH "$HOME/.local/share/fnm"
if [ -d "$FNM_PATH" ]
  fish_add_path "$FNM_PATH"
  fnm env --use-on-cd --shell fish | source
end
