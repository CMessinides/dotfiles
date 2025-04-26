if command -q uv
    uv --generate-shell-completion fish | source
end

if command -q uvx
    uvx --generate-shell-completion fish | source
end
