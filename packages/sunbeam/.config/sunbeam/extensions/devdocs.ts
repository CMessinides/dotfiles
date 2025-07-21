#!/usr/bin/env -S deno --allow-net=devdocs.io
import * as sunbeam from "sunbeam";

interface DevdocsDocset {
    name: string;
    slug: string;
    release?: string;
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
    ],
} as const satisfies sunbeam.Manifest;

if (Deno.args.length == 0) {
    console.log(JSON.stringify(manifest));
    Deno.exit(0);
}

const payload: sunbeam.Payload<typeof manifest> = JSON.parse(Deno.args[0]);

if (payload.command === "search-docsets") {
    const res = await fetch(`https://devdocs.io/docs/docs.json`);
    const docs: DevdocsDocset[] = await res.json();
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
} else {
    console.error("unknown command");
    Deno.exit(1);
}
