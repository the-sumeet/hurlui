<script lang="ts">
    import ace from "ace-builds";
    import "ace-builds/src-noconflict/theme-chaos"; // Example theme
    import "ace-builds/src-noconflict/mode-markdown"; // Markdown syntax support
    import { onMount, onDestroy } from "svelte";

    let { content = $bindable() } = $props();

    let theme = "chaos";
    let fontSize = 14;

    let editor: ace.Editor;
    let editorElement: HTMLElement;

    $effect(() => {
        if (editor) {
            editor.setValue(content, -1);
        }
    });

    // Initialize editor
    onMount(async () => {
        editor = ace.edit(editorElement, {
            // mode: `ace/mode/${mode}`,
            theme: `ace/theme/${theme}`,
            fontSize: `${fontSize}px`,
            // value: selectedFileContent,
        });

        // Sync with parent component
        editor.session.on("change", () => {
            content = editor.getValue();
        });

        // Add copy/paste keyboard shortcuts
        editor.commands.addCommand({
            name: "copy",
            bindKey: { win: "Ctrl-C", mac: "Cmd-C" },
            exec: function (editor) {
                const selectedText = editor.getSelectedText();
                if (selectedText) {
                    navigator.clipboard.writeText(selectedText);
                } else {
                    const currentLine = editor.session.getLine(
                        editor.getCursorPosition().row,
                    );
                    navigator.clipboard.writeText(currentLine);
                }
            },
        });

        editor.commands.addCommand({
            name: "paste",
            bindKey: { win: "Ctrl-V", mac: "Cmd-V" },
            exec: function (editor) {
                navigator.clipboard
                    .readText()
                    .then((text) => {
                        editor.insert(text);
                    })
                    .catch((err) => {
                        console.error("Failed to read clipboard:", err);
                    });
                return true;
            },
        });

        // Handle window resize
        const resizeObserver = new ResizeObserver(() => editor.resize());
        resizeObserver.observe(editorElement);

        onDestroy(() => {
            resizeObserver.disconnect();
            editor.destroy();
        });
    });

    $effect(() => {
        if (editor && content !== editor.getValue()) {
            editor.setValue(content, -1);
        }
    });
</script>

<div bind:this={editorElement} class="flex-1 ace-editor h-full"></div>
