vim.filetype.add({
  extension = {
    svx = "markdown",
    mdx = "markdown",
    svelte = "svelte",
    patch = "patch",
  },
  filename = {
    [".prettierrc"] = "jsonc",
    [".eslintrc"] = "jsonc",
    ["tsconfig.json"] = "jsonc",
    ["jsconfig.json"] = "jsonc",
  },
})
