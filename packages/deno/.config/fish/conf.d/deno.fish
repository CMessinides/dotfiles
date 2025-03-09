set DENO_INSTALL "$HOME/.deno"
if [ -d "$DENO_INSTALL" ]
  fish_add_path "$DENO_INSTALL"
end

source "/home/cam/.deno/env.fish"