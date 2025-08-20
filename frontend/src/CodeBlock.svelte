<script lang="ts">
    import ace from "ace-builds";
    import "ace-builds/src-noconflict/theme-chaos"; // Example theme
    import "ace-builds/src-noconflict/mode-html";
    import "ace-builds/src-noconflict/mode-json";
    import "ace-builds/src-noconflict/mode-xml";
    import { onMount, onDestroy } from "svelte";
    import { getContext } from "svelte";
    import { appState } from "./state.svelte";

    let { value } = $props();

    let theme = "chaos";
    let mode = "markdown";
    let fontSize = 14;

    let editor: ace.Editor;
    let editorElement: HTMLElement;

    // $effect(() => {
    //     if (appState.bold && editor) {
    //         bold();
    //         appState.bold = false;
    //     } else if (appState.italic && editor) {
    //         italic();
    //         appState.italic = false;
    //     } else if (appState.underline && editor) {
    //         underline();
    //         appState.underline = false;
    //     }
    // });

    // $effect(() => {
    //     if (appState.selectedNote != null) {
    //         if (editor) {
    //             editor.setValue(appState.selectedNote.Content, -1);
    //         }
    //     } else {
    //         if (editor) {
    //             editor.setValue("", -1);
    //         }
    //     }
    // });

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
