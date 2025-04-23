-- [[ Basic Keymaps ]]

-- Keymaps for better default experience
-- See `:help vim.keymap.set()`
vim.keymap.set({ "n", "v" }, "<Space>", "<Nop>", { silent = true })

-- Keep cursor in the middle of the screen when jumping by half-pages
vim.keymap.set("n", "<C-d>", "<C-d>zz")
vim.keymap.set("n", "<C-u>", "<C-u>zz")

-- Put the highlighted search term in the middle of the screen
-- when jumping between terms.
vim.keymap.set("n", "n", "nzzzv")
vim.keymap.set("n", "N", "Nzzzv")

-- Remap for dealing with word wrap
vim.keymap.set("n", "k", "v:count == 0 ? 'gk' : 'k'", { expr = true, silent = true })
vim.keymap.set("n", "j", "v:count == 0 ? 'gj' : 'j'", { expr = true, silent = true })

-- Remap the default tab* keymaps (conflict with iTerm2)
vim.keymap.set("n", "<Leader>tc", vim.cmd.tabnew, { desc = "[C]reate [T]ab" })
vim.keymap.set("n", "<Leader>tn", vim.cmd.tabnext, { desc = "[N]ext [T]ab" })
vim.keymap.set("n", "<Leader>tp", vim.cmd.tabprev, { desc = "[P]rev [T]ab" })

-- Add keymaps for buffer commands
vim.keymap.set("n", "<Leader>bd", vim.cmd.bd, { desc = "[B]uffer [D]elete" })
vim.keymap.set("n", "<Leader>bda", "<Cmd>%bd<CR>", { desc = "[B]uffer [D]elete [A]ll" })
