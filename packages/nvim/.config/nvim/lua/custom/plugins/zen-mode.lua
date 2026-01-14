return {
    {
        "folke/zen-mode.nvim",
        cmd = "ZenMode",
        opts = {
            window = {
                width = 80,
            },
            plugins = {
                kitty = {
                    enabled = true,
                },
            },
        },
        keys = {
            { "<leader>uz", "<cmd>ZenMode<cr>", desc = "Toggle Zen mode" },
        },
    },
}
