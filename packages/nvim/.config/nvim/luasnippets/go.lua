local line_begin = require("luasnip.extras.expand_conditions").line_begin

return {
    s(
        {
            trig = "ierr",
            dscr = "Implement the error interface for a type",
            condition = line_begin,
        },
        fmta(
            [[
                func (<> <>) Error() string {
                    <>
                }
            ]],
            {
                i(1, "e"),
                i(2, "*ErrorType"),
                i(0),
            }
        )
    ),
    s(
        {
            trig = "cerr",
            dscr = "Implement a custom error type",
            condition = line_begin,
        },
        fmta(
            [[
                type <> struct {
                    <>
                    Err error
                }

                func (e *<>) Error() string {
                    return "<>: " + e.Err.Error()
                }

                func (e *<>) Unwrap() error {
                    return e.Err
                }
            ]],
            {
                i(1, "CustomError"),
                i(2),
                rep(1),
                i(3),
                rep(1),
            }
        )
    ),
}
