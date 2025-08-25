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
        keys = {
            {
                "<leader>tt",
                function()
                    require("neotest").summary.toggle()
                end,
                desc = "[T]ests: [T]oggle Summary",
            },
            {
                "<leader>tr",
                function()
                    require("neotest").run.run()
                end,
                desc = "[T]ests: [R]un Nearest",
            },
            {
                "<leader>ts",
                function()
                    require("neotest").run.run({ suite = true })
                end,
                desc = "[T]ests: Run [S]uite",
            },
            {
                "<leader>tsw",
                function()
                    require("neotest").watch.toggle({ suite = true })
                end,
                desc = "[T]ests: Toggle [S]uite [W]atcher",
            },
            {
                "<leader>tf",
                function()
                    require("neotest").run.run(vim.fn.expand("%"))
                end,
                desc = "[T]ests: Run [F]ile",
            },
            {
                "<leader>tfw",
                function()
                    require("neotest").watch.toggle(vim.fn.expand("%"))
                end,
                desc = "[T]ests: Toggle [F]ile [W]atcher",
            },
            {
                "<leader>to",
                function()
                    require("neotest").output.open()
                end,
                desc = "[T]ests: Open [O]utput",
            },
            {
                "[t",
                function()
                    require("neotest").jump.prev({ status = "failed" })
                end,
                desc = "Tests: Previous failed test",
                silent = true,
            },
            {
                "]t",
                function()
                    require("neotest").jump.next({ status = "failed" })
                end,
                desc = "Tests: Next failed test",
                silent = true,
            },
            {
                "[T",
                function()
                    require("neotest").jump.prev()
                end,
                desc = "Tests: Previous test",
                silent = true,
            },
            {
                "]T",
                function()
                    require("neotest").jump.next()
                end,
                desc = "Tests: Next test",
                silent = true,
            },
        },
        config = function()
            require("neotest").setup({
                adapters = {
                    require("neotest-vitest"),
                },
            })
        end,
    },
}
