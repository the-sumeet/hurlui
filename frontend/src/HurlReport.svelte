<script lang="ts">
    import type { main } from "wailsjs/go/models";
    import * as Tooltip from "$lib/components/ui/tooltip/index.js";
    import { Info } from "lucide-svelte";
    import * as Select from "$lib/components/ui/select/index.js";

    import ResponseCall from "./ResponseCall.svelte";
    let {
        hurlReport,
    }: {
        hurlReport: main.HurlSession[] | null;
    } = $props();

    // let entries: main.HurlEntry[] = $state([]);
    let selectedSessionIndex: string = $state("0");
    // $effect(() => {
    //     const index = parseInt(selectedSessionIndex);
    //     triggerContent = hurlReport !== null ? hurlReport[index].filename : "";
    //     entries = hurlReport !== null ? hurlReport[index].entries : [];
    // });
    let triggerContent = $derived.by(() => {
        if (!hurlReport || hurlReport.length === 0) return "";
        const index = Number.parseInt(selectedSessionIndex) || 0;
        const safeIndex = Math.min(Math.max(index, 0), hurlReport.length - 1);
        return hurlReport[safeIndex]?.filename ?? "";
    });

    let entries = $derived.by(() => {
        if (!hurlReport || hurlReport.length === 0)
            return [] as main.HurlEntry[];
        const index = Number.parseInt(selectedSessionIndex) || 0;
        const safeIndex = Math.min(Math.max(index, 0), hurlReport.length - 1);
        return hurlReport[safeIndex]?.entries ?? ([] as main.HurlEntry[]);
    });

    // Ensure selected index remains valid when hurlReport updates
    $effect(() => {
        if (!hurlReport || hurlReport.length === 0) return;
        const idx = Number.parseInt(selectedSessionIndex) || 0;
        if (idx < 0 || idx >= hurlReport.length) selectedSessionIndex = "0";
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
</script>

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
            {#each entries as entry (entry.index)}
                {#each entry.calls as call, j (`${call.request.url}:${call.response.status ?? 0}:${j}`)}
                    <ResponseCall
                        showCallNumber={entry.calls.length > 1}
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
