return {
    {
        'stevearc/oil.nvim',
        -- Optional dependencies
        dependencies = { "nvim-tree/nvim-web-devicons" },
        config = function()
            require('oil').setup({
                view_options = {
                    show_hidden = true,
                    is_always_hidden = function(name)
                        return name == ".." or name == ".git"
                    end
                },
            })

            vim.keymap.set('n', '<leader>.', '<cmd>Oil<cr>', { silent = true, desc = "Open parent directory" })
            vim.keymap.set(
                'n',
                '<leader>onv',
                function() require('oil').open(vim.fn.stdpath('config')) end,
                { silent = true, desc = "Open neovim config directory" }
            )
            vim.keymap.set(
                'n',
                '<leader>odf',
                function()
                    local dotfiles = vim.env.DOTFILES
                    if dotfiles then
                        require('oil').open(dotfiles)
                    end
                end,
                { silent = true, desc = "Open dotfiles directory" }
            )
        end
    }
}
