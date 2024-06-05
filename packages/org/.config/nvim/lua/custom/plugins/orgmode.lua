return {
    {
        'nvim-orgmode/orgmode',
        cmd = {
            'OrgCapture'
        },
        ft = 'org',
        event = 'VeryLazy',
        config = function()
            -- Setup orgmode
            require('orgmode').setup({
                org_agenda_files = { '~/org/**/*', '~/private/**/*' },
                org_default_notes_file = '~/org/refile.org',
            })

            vim.api.nvim_create_user_command('OrgCapture', function()
                require('orgmode').capture:prompt()
            end, { desc = "Open org capture prompt" })
        end,
    }
}
