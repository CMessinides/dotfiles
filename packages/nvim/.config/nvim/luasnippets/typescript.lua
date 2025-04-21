return {
    s(
        {
            trig = "cenum",
            dscr = "Create an enum type with a const object"
        },
        fmta(
            [[
                    export const <> = {
                        <>
                    } as const

                    export type <> = (typeof <>)[keyof typeof <>]
                ]],
            {
                i(1, ""),
                i(0),
                rep(1),
                rep(1),
                rep(1)
            }
        )
    ),
}
