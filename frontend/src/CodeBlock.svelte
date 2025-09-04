<script lang="ts">
    import ace from "ace-builds";
    import "ace-builds/src-noconflict/theme-chaos"; // Example theme
    import "ace-builds/src-noconflict/mode-html";
    import "ace-builds/src-noconflict/mode-json";
    import "ace-builds/src-noconflict/mode-xml";
    import { onMount, onDestroy } from "svelte";

    let { value, mode } = $props();

    $effect(() => {
        if (mode && editor) {
            const session = editor.getSession();
            session.setMode(`ace/mode/${mode}`);
        }
    });

    // Update editor content when value prop changes
    $effect(() => {
        if (!editor) return;
        const session = editor.getSession();
        const current = session.getValue();
        if (current !== value) {
            editor.setValue(value ?? "", -1);
        }
    });

    let theme = "chaos";
    let fontSize = 14;

    let editor: ace.Editor;
    let editorElement: HTMLElement;

    function formatJson(input: string): string {
        try {
            const parsed = JSON.parse(input);
            return JSON.stringify(parsed, null, 2);
        } catch (_) {
            return input;
        }
    }

    function formatXml(input: string): string {
        if (!input) return input;
        const reg = /(>)(<)(\/*)/g;
        let xml = input.replace(reg, "$1\n$2$3");
        const PADDING = "  ";
        let formatted = "";
        let pad = 0;
        for (const node of xml.split("\n")) {
            if (!node) continue;
            let indentChange = 0;
            if (/^<\/.+>/.test(node)) {
                pad = Math.max(pad - 1, 0);
            } else if (/^<[^!?][^>]*[^\/]>(?!.*<\/)/.test(node)) {
                indentChange = 1;
            } else if (/^<[^!?].*\/>/.test(node)) {
                indentChange = 0;
            } else if (/.+<\/[^>]+>$/.test(node)) {
                indentChange = 0;
            }
            formatted += PADDING.repeat(pad) + node + "\n";
            pad += indentChange;
        }
        return formatted.trim();
    }

    // Public method: pretty-format current editor content based on `mode`
    export function format() {
        const current = (editor?.getValue?.() ?? value ?? "") as string;
        if (!current) return;
        let next = current;
        if (mode === "json") {
            next = formatJson(current);
        } else if (mode === "xml" || mode === "html") {
            next = formatXml(current);
        }
        if (next !== current) {
            editor?.setValue?.(next, -1);
        }
    }

    // Initialize editor
    onMount(async () => {
        editor = ace.edit(editorElement, {
            mode: `ace/mode/${mode}`,
            theme: `ace/theme/${theme}`,
            fontSize: `${fontSize}px`,
            // value: selectedFileContent,
            value: value,
        });

        // Propagate editor changes back to bound `value`
        editor.on("change", () => {
            // Avoid unnecessary churn: only update if changed
            const newVal = editor.getValue() ?? "";
            if (newVal !== value) value = newVal;
        });

        // editor.setValue(value, -1);

        // editor.commands.addCommands([
        //     {
        //         name: "find",
        //         bindKey: { win: "Ctrl-F", mac: "Cmd-F" },
        //         exec: function (editor) {
        //             editor.execCommand("find");
        //             editor.searchBox.show();
        //         },
        //     },
        // ]);

        // Handle window resize
        const resizeObserver = new ResizeObserver(() => editor.resize());
        resizeObserver.observe(editorElement);

        onDestroy(() => {
            resizeObserver.disconnect();
            editor.destroy();
        });
    });
</script>

<div
    bind:this={editorElement}
    class="flex-1 ace-editor h-full rounded-xl"
></div>
