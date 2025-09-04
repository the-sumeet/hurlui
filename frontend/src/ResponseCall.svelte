<script lang="ts">
    import * as Card from "$lib/components/ui/card/index.js";
    import * as Tabs from "$lib/components/ui/tabs/index.js";
    import * as Select from "$lib/components/ui/select/index.js";
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { main } from "../wailsjs/go/models";
    import CodeBlock from "./CodeBlock.svelte";
    import { Textarea } from "$lib/components/ui/textarea/index.js";
    import { snakeToTitleCase } from "./utils";
    import { responseTypes, timingsDescription } from "./constants";
    import * as Tooltip from "$lib/components/ui/tooltip/index.js";
    import { Info } from "lucide-svelte";
    import { Button } from "$lib/components/ui/button/index.js";
    import type { SvelteComponent } from "svelte";
    import { Wand } from "lucide-svelte";
    let {
        showCallNumber,
        callNumber,
        call,
    }: { showCallNumber: boolean; callNumber: number; call: main.HurlCall } =
        $props();

    let codeBlockRef: SvelteComponent;

    function getResponseHeaders(call: main.HurlCall): string {
        let result = "";
        call.response.headers.forEach((element) => {
            result += `${element.name}: ${element.value}\n`;
        });
        return result;
    }

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

    let responseType = $state(getResponseMode(call));

    // Local editable/format-able body content used by the CodeBlock
    let responseBody: string = $state(
        (call.response.bodyContent || call.response.body) as string,
    );

    const statusCode = call.response.status;

    const statusClass = (function () {
        const sc = statusCode ?? 0;
        if (sc >= 200 && sc < 300) return "text-green-600";
        if (sc >= 300 && sc < 400) return "text-yellow-600";
        if (sc >= 400 && sc < 500) return "text-orange-600";
        if (sc >= 500) return "text-red-600";
        return "text-muted-foreground";
    })();

    const triggerContent = $derived(
        responseTypes.find((f) => f.value === responseType)?.label ??
            "Select a type",
    );
</script>

<div class="flex flex-col only:h-full h-[600px]">
    {#if showCallNumber}
        <p class="text-sm">Call #{callNumber}</p>
    {/if}

    <Card.Root class="h-full py-1 gap-1">
        <Card.Header class="px-1">
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
        <Card.Content class="h-full px-1">
            <!-- Main tab -->
            <Tabs.Root
                value="response"
                class="h-full w-full border rounded-xl shadow-sm p-1"
            >
                <Tabs.List>
                    <Tabs.Trigger value="response">Response</Tabs.Trigger>
                    <Tabs.Trigger value="request">Request</Tabs.Trigger>
                    <Tabs.Trigger value="timing">Timing</Tabs.Trigger>
                </Tabs.List>
                <!-- Response Content -->
                <Tabs.Content value="response">
                    <Tabs.Root
                        value="body"
                        class="h-full w-full border rounded-xl shadow-sm p-1"
                    >
                        <Tabs.List>
                            <Tabs.Trigger value="body">Body</Tabs.Trigger>
                            <Tabs.Trigger value="password"
                                >Response</Tabs.Trigger
                            >
                        </Tabs.List>
                        <Tabs.Content value="body">
                            <div class="flex flex-col h-full gap-1">
                                <div class="flex items-center gap-2">
                                    <Select.Root
                                        type="single"
                                        name="responseType"
                                        bind:value={responseType}
                                    >
                                        <Select.Trigger class="w-[180px]">
                                            {triggerContent}
                                        </Select.Trigger>
                                        <Select.Content>
                                            <Select.Group>
                                                <Select.Label
                                                    >Response Type</Select.Label
                                                >
                                                {#each responseTypes as responseType (responseType.value)}
                                                    <Select.Item
                                                        value={responseType.value}
                                                        label={responseType.label}
                                                    >
                                                        {responseType.label}
                                                    </Select.Item>
                                                {/each}
                                            </Select.Group>
                                        </Select.Content>
                                    </Select.Root>
                                    <Button
                                        variant="secondary"
                                        onclick={() => codeBlockRef?.format()}
                                    >
                                        <Wand />
                                    </Button>
                                </div>
                                <div class="flex-1 h-full border rounded-xl">
                                    <CodeBlock
                                        bind:this={codeBlockRef}
                                        value={call.response.bodyContent ||
                                            call.response.body}
                                        mode={responseType}
                                    />
                                </div>
                            </div>
                        </Tabs.Content>
                        <Tabs.Content value="password"
                            ><CodeBlock
                                value={getResponseHeaders(call)}
                                mode="json"
                            /></Tabs.Content
                        >
                    </Tabs.Root>
                </Tabs.Content>
                <!-- Request Content -->
                <Tabs.Content value="request"
                    ><Textarea
                        class="h-full"
                        readonly
                        value={getEntryText(call)}
                    /></Tabs.Content
                >
                <!-- Timing Content -->
                <Tabs.Content value="timing">
                    <div class="flex flex-col">
                        {#each Object.entries(call.timings) as [key, value]}
                            {@render timingRow(
                                snakeToTitleCase(key) || key,
                                value,
                                timingsDescription[key] || null,
                            )}
                        {/each}
                    </div>
                </Tabs.Content>
            </Tabs.Root>
        </Card.Content>
        <Card.Footer>
            <p class={statusClass}>{statusCode ?? ""}</p>
        </Card.Footer>
    </Card.Root>
</div>

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
