return {
    {
        "folke/twilight.nvim",
        cmd = {
            "Twilight",
            "TwilightEnable",
            "TwilightDisable"
        },
        keys = {
            { "<leader>tt", "<cmd>Twilight<cr>", desc = "Toggle Twilight" }
        },
        opts = {
            context = 5,
            expand = {
                "function",
                "method",
                "table",
                "if_statement",
                -- org mode
                "paragraph",
                "headline",
                "list",
            },
        }
    }
}
