return {
    {
        "nvim-neotest/neotest",
        dependencies = {
            "nvim-neotest/nvim-nio",
            "nvim-lua/plenary.nvim",
            "antoinemadec/FixCursorHold.nvim",
            "nvim-treesitter/nvim-treesitter",

            -- Adapters
            "marilari88/neotest-vitest",
        },
        config = function()
            require("neotest").setup({
                adapters = {
                    require("neotest-vitest"),
                },
            })

            vim.api.nvim_create_autocmd("BufRead", {
                pattern = { "*.test.{ts,js,tsx,jsx}" },
                callback = function()
                    vim.keymap.set("n", "<leader>nw", function()
                        local neotest = require("neotest")
                        local filename = vim.fn.expand("%")
                        if neotest.watch.is_watching(filename) then
                            neotest.watch.stop(filename)
                        else
                            neotest.watch.watch(filename)
                        end
                    end, {
                        desc = "[N]eotest: Toggle [W]atch File",
                        buffer = true,
                    })
                end,
            })
        end,
    },
}
