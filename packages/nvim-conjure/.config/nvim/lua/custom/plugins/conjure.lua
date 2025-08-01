return {
    -- Interactive REPL editing
    {
        "Olical/conjure",
        cmd = "ConjureSchool",
        ft = { "clojure", "fennel", "python" }, -- etc
        -- [Optional] cmp-conjure for cmp
        dependencies = {
            {
                "PaterJason/cmp-conjure",
                config = function()
                    local cmp = require("cmp")
                    local config = cmp.get_config()
                    table.insert(config.sources, {
                        name = "buffer",
                        option = {
                            sources = {
                                { name = "conjure" },
                            },
                        },
                    })
                    cmp.setup(config)
                end,
            },
        },
        config = function()
            require("conjure.main").main()
            require("conjure.mapping")["on-filetype"]()
        end,
        init = function()
            -- Set configuration options here
            -- vim.g["conjure#debug"] = true

            -- Disable the documentation mapping
            vim.g["conjure#mapping#doc_word"] = false
        end,
    },

    -- Structural editing, optional
    "guns/vim-sexp",
    "tpope/vim-sexp-mappings-for-regular-people",
}
