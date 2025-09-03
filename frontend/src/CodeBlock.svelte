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

    // Initialize editor
    onMount(async () => {
        editor = ace.edit(editorElement, {
            mode: `ace/mode/${mode}`,
            theme: `ace/theme/${theme}`,
            fontSize: `${fontSize}px`,
            // value: selectedFileContent,
            value: value,
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
