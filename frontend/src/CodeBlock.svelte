<script lang="ts">
    import ace from "ace-builds";
    import "ace-builds/src-noconflict/theme-chaos"; // Example theme
    import "ace-builds/src-noconflict/mode-html";
    import "ace-builds/src-noconflict/mode-json";
    import "ace-builds/src-noconflict/mode-xml";
    import { onMount, onDestroy } from "svelte";
    import { getContext } from "svelte";
    import { appState } from "./state.svelte";

    let { value, mode } = $props();

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
        });

        editor.setValue(value, -1);

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
