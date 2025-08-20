<script lang="ts">
    import type { main } from "wailsjs/go/models";
    import { Textarea } from "$lib/Components/ui/textarea/index.js";
    import * as Tooltip from "$lib/components/ui/tooltip/index.js";
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { Info } from "lucide-svelte";
    import { timingsDescription } from "./constants";
    import { snakeToTitleCase } from "./utils";
    import CodeBlock from "./CodeBlock.svelte";
    import * as Select from "$lib/Components/ui/select/index.js";

    import ResponseCall from "./ResponseCall.svelte";
    let {
        hurlReport,
    }: {
        hurlReport: main.HurlSession[] | null;
    } = $props();

    let triggerContent = $state("");
    let entries: main.HurlEntry[] = $state([]);
    let selectedSessionIndex: string = $state("0");
    $effect(() => {
        const index = parseInt(selectedSessionIndex);
        triggerContent = hurlReport !== null ? hurlReport[index].filename : "";
        entries = hurlReport !== null ? hurlReport[index].entries : [];
    });

    function getEntryText(call: main.HurlCall): string {
        let result = "";

        const method = call.request?.method || "";
        const url = call.request?.url || "";
        const query_string = call.request?.query_string || "";
        let headers = "";
        call?.request.headers.forEach((element) => {
            headers += `${element.name}: ${element.value}\n`;
        });

        result += `${method} ${url}?${query_string}\n${headers}\n`;

        return result;
    }

    function getResponseHeaders(call: main.HurlCall): string {
        let result = "";
        call.response.headers.forEach((element) => {
            result += `${element.name}: ${element.value}\n`;
        });
        return result;
    }

    function getResponseMode(call: main.HurlCall): string {
        const contentTypeHeader = call.response.headers.find(
            (header) => header.name.toLowerCase() === "content-type",
        );

        if (!contentTypeHeader) return "html";

        const contentType = contentTypeHeader.value.toLowerCase();

        if (
            contentType.includes("application/json") ||
            contentType.includes("text/json")
        ) {
            return "json";
        } else if (
            contentType.includes("application/xml") ||
            contentType.includes("text/xml")
        ) {
            return "xml";
        } else {
            return "html";
        }
    }
</script>

{#snippet timingRow(key: string, value: string, tooltip: string | null)}
    <div class="flex gap-1">
        {#if tooltip}
            <Tooltip.Provider>
                <Tooltip.Root>
                    <Tooltip.Trigger><Info size={16} /></Tooltip.Trigger>
                    <Tooltip.Content>
                        <p>{tooltip}</p>
                    </Tooltip.Content>
                </Tooltip.Root>
            </Tooltip.Provider>
        {/if}
        <p class="text-sm">{key}: {value}</p>
    </div>
{/snippet}

<div class="flex-1 flex flex-col overflow-y-hidden p-1 h-full">
    <!-- Files/sessions list -->
    {#if hurlReport && hurlReport.length > 0}
        <!-- Files/sessions list -->
        {#if hurlReport && hurlReport.length > 1}
            <Select.Root type="single" bind:value={selectedSessionIndex}>
                <Select.Trigger class="w-full">{triggerContent}</Select.Trigger>
                <Select.Content>
                    {#each hurlReport || [] as session, index}
                        <Select.Item value={index.toString()}>
                            {session.filename}
                        </Select.Item>
                    {/each}
                </Select.Content>
            </Select.Root>
        {/if}

        <!-- Response container -->
        <div class="flex-1 flex flex-col gap-1 overflow-y-scroll h-full">
            {#each entries as trigger, i}
                {#each trigger.calls as call, j}
                    <ResponseCall
                        showCallNumber={trigger.calls.length > 1}
                        callNumber={j + 1}
                        {call}
                    />
                {/each}
            {/each}
        </div>
    {:else}
        <div
            class="flex flex-col items-center justify-center h-full text-gray-600 italic"
        >
            <p>Send some requests</p>
        </div>{/if}
</div>
