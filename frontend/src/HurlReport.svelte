<script lang="ts">
    import * as Card from "$lib/Components/ui/card/index.js";
    import { Button } from "$lib/Components/ui/button/index.js";
    import * as Tabs from "$lib/Components/ui/tabs/index.js";
    import * as Select from "$lib/Components/ui/select/index.js";
    import type { main } from "wailsjs/go/models";
    import { Textarea } from "$lib/Components/ui/textarea/index.js";
    import * as Tooltip from "$lib/components/ui/tooltip/index.js";
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { Info } from "lucide-svelte";
    import { timingsDescription } from "./constants";
    import { snakeToTitleCase } from "./utils";
    import CodeBlock from "./CodeBlock.svelte";
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

<div class="flex-1 flex flex-col h-full overflow-y-hide p-1">
    <!-- Files/sessions list -->
    {#if hurlReport && hurlReport.length > 0}
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

        <div class="flex flex-col gap-1 overflow-y-auto h-full mt-1">
            {#each entries as trigger, i}
                {#each trigger.calls as call, j}
                    <div
                        class="flex flex-col {trigger.calls.length > 1
                            ? 'border rounded-xl shadow-sm'
                            : ''} gap-1 p-1"
                    >
                        {#if trigger.calls.length > 1}
                            <p>Call #{j + 1}</p>
                        {/if}

                        <Card.Root class="py-2 gap-2">
                            <Card.Header class="px-2">
                                <!-- <Card.Title>{trigger}</Card.Title> -->
                                <Card.Description>
                                    <div class="flex gap-1 items-start">
                                        <Badge variant="outline" class="h-min"
                                            >{call.request.method}</Badge
                                        >
                                        <div class="flex-1">
                                            {call.request.url}
                                        </div>
                                    </div></Card.Description
                                >
                                <!-- <Card.Action>
                                    <Button variant="link" size="sm"
                                        >Load Response</Button
                                    >
                                </Card.Action> -->
                            </Card.Header>
                            <Card.Content class="px-2">
                                <!-- Main tab -->
                                <Tabs.Root
                                    value="response"
                                    class="w-full border rounded-xl shadow-sm p-1"
                                >
                                    <Tabs.List>
                                        <Tabs.Trigger value="response"
                                            >Response</Tabs.Trigger
                                        >
                                        <Tabs.Trigger value="request"
                                            >Request</Tabs.Trigger
                                        >
                                        <Tabs.Trigger value="timing"
                                            >Timing</Tabs.Trigger
                                        >
                                    </Tabs.List>
                                    <!-- Response Content -->
                                    <Tabs.Content value="response">
                                        <Tabs.Root
                                            value="body"
                                            class="w-full border rounded-xl shadow-sm p-1"
                                        >
                                            <Tabs.List>
                                                <Tabs.Trigger value="body"
                                                    >Body</Tabs.Trigger
                                                >
                                                <Tabs.Trigger value="password"
                                                    >Response</Tabs.Trigger
                                                >
                                            </Tabs.List>
                                            <Tabs.Content value="body">
                                                <div
                                                    class="h-48 border rounded-xl"
                                                >
                                                    <CodeBlock
                                                        value={call.response
                                                            .bodyContent}
                                                    />
                                                </div>
                                            </Tabs.Content>
                                            <Tabs.Content value="password"
                                                >Change your password here.</Tabs.Content
                                            >
                                        </Tabs.Root>
                                    </Tabs.Content>
                                    <!-- Request Content -->
                                    <Tabs.Content value="request"
                                        ><Textarea
                                            readonly
                                            value={getEntryText(call)}
                                            rows={6}
                                        /></Tabs.Content
                                    >
                                    <!-- Timing Content -->
                                    <Tabs.Content value="timing">
                                        <div class="flex flex-col">
                                            {#each Object.entries(call.timings) as [key, value]}
                                                {@render timingRow(
                                                    snakeToTitleCase(key) ||
                                                        key,
                                                    value,
                                                    timingsDescription[key] ||
                                                        null,
                                                )}
                                            {/each}
                                        </div>
                                    </Tabs.Content>
                                </Tabs.Root>
                            </Card.Content>
                            <Card.Footer>
                                <p class="text-green-600">200 OK</p>
                            </Card.Footer>
                        </Card.Root>
                    </div>
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
