#!/usr/bin/env -S deno --allow-net=devdocs.io
import * as sunbeam from "../lib/sunbeam/mod.ts";

interface DevDocsDocset {
    name: string;
    slug: string;
    release?: string;
}

interface DevDocsEntry {
    name: string;
    type: string;
    path: string;
}

interface DevDocsEntryManifest {
    entries: DevDocsEntry[];
}

const manifest = {
    title: "Devdocs",
    description: "Search the devdocs.io documentation",
    commands: [
        {
            name: "search-docsets",
            title: "Search docsets",
            mode: "filter",
        },
        {
            name: "search-entries",
            title: "Search entries",
            mode: "filter",
            params: [
                {
                    name: "docset",
                    type: "string",
                    title: "Docset Slug",
                },
            ],
        },
    ],
} as const satisfies sunbeam.Manifest;

if (Deno.args.length == 0) {
    console.log(JSON.stringify(manifest));
    Deno.exit(0);
}

const payload: sunbeam.Payload<typeof manifest> = JSON.parse(Deno.args[0]);

if (payload.command === "search-docsets") {
    const res = await fetch(`https://devdocs.io/docs/docs.json`);
    const docs: DevDocsDocset[] = await res.json();
    const list: sunbeam.List = {
        items: docs.map((doc) => ({
            title: doc.name,
            subtitle: doc.release || "latest",
            accessories: [doc.slug],
            actions: [
                {
                    type: "run",
                    title: `Search ${doc.name} entries`,
                    command: "search-entries",
                    params: {
                        docset: doc.slug,
                    },
                },
                {
                    type: "open",
                    title: "Open in browser",
                    url: `https://devdocs.io/${doc.slug}`,
                    exit: true,
                },
            ],
        })),
    };

    console.log(JSON.stringify(list));
} else if (payload.command === "search-entries") {
    const { docset } = payload.params;
    const res = await fetch(`https://devdocs.io/docs/${docset}/index.json`);
    const { entries }: DevDocsEntryManifest = await res.json();
    const list: sunbeam.List = {
        items: entries.map<sunbeam.ListItem>((entry) => {
            const url = `https://devdocs.io/${docset}/${entry.path}`;
            return {
                title: entry.name,
                subtitle: entry.type,
                actions: [
                    {
                        title: "Open in Browser",
                        type: "open",
                        url,
                    },
                    {
                        title: "Copy URL",
                        type: "copy",
                        key: "c",
                        text: url,
                    },
                ],
            };
        }),
    };

    console.log(JSON.stringify(list));
} else {
    console.error("unknown command");
    Deno.exit(1);
}
