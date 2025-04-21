return {
    s(
        {
            trig = 'sstory',
            dscr = 'Setup Storybook story'
        },
        fmt(
            [[
                import { Meta, StoryObj } from "@storybook/react";
                import { (1) } from "./(2)";

                const meta: Meta<typeof (2)> = {
                    title: "Components/(2)",
                    component: (2),
                    args: {},
                    argTypes: {},
                    tags: ["autodocs"],
                };

                export default meta;

                type Story = StoryObj<typeof (2)>;

                export const Default: Story = {};
            ]],
            {
                i(1, 'ComponentName'),
                rep(1),
            },
            {
                delimiters = '()',
            }
        )
    ),
}
