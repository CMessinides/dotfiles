-- Telescope integration adapted from https://github.com/joaomsa/telescope-orgmode.nvim
-- MIT License, Copyright (c) 2022 Joao Sa

local get_entries = function(opts)
    local orgmode = require('orgmode.api')
    vim.print(orgmode)
    local file_results = vim.tbl_map(function(file)
        return { file = file, filename = file.filename }
    end, orgmode.load())

    if not opts.archived then
        file_results = vim.tbl_filter(function(entry)
            return not entry.file.is_archive_file
        end, file_results)
    end

    if opts.max_depth == 0 then
        return file_results
    end

    local results = {}
    for _, file_entry in ipairs(file_results) do
        for _, headline in ipairs(file_entry.file.headlines) do
            local allowed_depth = opts.max_depth == nil or headline.level <= opts.max_depth
            local allowed_archive = opts.archived or not headline.is_archived
            if allowed_depth and allowed_archive then
                local entry = {
                    file = file_entry.file,
                    filename = file_entry.filename,
                    headline = headline
                }
                table.insert(results, entry)
            end
        end
    end

    return results
end

local make_entry = function(opts)
    local entry_display = require("telescope.pickers.entry_display")

    local displayer = entry_display.create({
        separator = ' ',
        items = {
            { width = vim.F.if_nil(opts.location_width, 20) },
            { remaining = true }
        }
    })

    local function make_display(entry)
        return displayer({ entry.location, entry.line })
    end

    return function(entry)
        local headline = entry.headline

        local lnum = nil
        local location = vim.fn.fnamemodify(entry.filename, ':t')
        local line = ""

        if headline then
            lnum = headline.position.start_line
            location = string.format('%s:%i', location, lnum)
            line = string.format('%s %s', string.rep('*', headline.level), headline.title)
        end

        return {
            value = entry,
            ordinal = location .. ' ' .. line,
            filename = entry.filename,
            lnum = lnum,
            display = make_display,
            location = location,
            line = line
        }
    end
end

local search_headings = function(opts)
    local pickers = require("telescope.pickers")
    local finders = require("telescope.finders")
    local conf = require("telescope.config").values
    opts = opts or {}

    pickers.new(opts, {
        prompt_title = "Search Headings",
        finder = finders.new_table {
            results = get_entries(opts),
            entry_maker = opts.entry_maker or make_entry(opts),
        },
        sorter = conf.generic_sorter(opts),
        previewer = conf.grep_previewer(opts),
    }):find()
end

return {
    {
        'nvim-orgmode/orgmode',
        cmd = {
            'OrgCapture'
        },
        ft = 'org',
        keys = {
            {
                '<leader>so',
                search_headings,
                desc = '[S]earch [O]rg headlines',
            }
        },
        event = 'VeryLazy',
        config = function()
            -- Setup orgmode
            require('orgmode').setup({
                org_agenda_files = { '~/org/**/*', '~/private/**/*' },
                org_default_notes_file = '~/org/refile.org',
                org_capture_templates = {
                    t = {
                        description = 'Task',
                        template = '* TODO %?\n  %u',
                    },
                    m = {
                        description = 'Meeting',
                        template = '\n**** %?\n     %^U',
                        datetree = {
                            tree_type = 'day',
                        },
                        target = '~/org/calendar.org',
                    },
                },
            })

            vim.api.nvim_create_user_command('OrgCapture', function()
                require('orgmode').capture:prompt()
            end, { desc = "Open org capture prompt" })

            require('which-key').register {
                ['<leader>o'] = { name = '[O]rg', _ = 'which_key_ignore' },
            }

            vim.api.nvim_create_autocmd('Filetype', {
                pattern = 'org',
                callback = function(event)
                    require('which-key').register({
                        ['<leader>ob'] = { name = '[O]rg [B]abel', _ = 'which_key_ignore' },
                        ['<leader>od'] = { name = '[O]rg [D]ate', _ = 'which_key_ignore' },
                        ['<leader>oi'] = { name = '[O]rg [I]nsert', _ = 'which_key_ignore' },
                        ['<leader>ol'] = { name = '[O]rg [L]ink', _ = 'which_key_ignore' },
                        ['<leader>ox'] = { name = '[O]rg Timesheet', _ = 'which_key_ignore' },
                    }, {
                        buffer = event.buf
                    })
                end
            })
        end,
    }
}
