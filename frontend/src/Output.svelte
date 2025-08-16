<script lang="ts">
    import ace from "ace-builds";
    import "ace-builds/src-noconflict/theme-github"; // Example theme
    import "ace-builds/src-noconflict/mode-markdown"; // Markdown syntax support
    import { onMount, onDestroy } from "svelte";
    import { getContext } from "svelte";
    import { appState } from "./state.svelte";

    let theme = "github";
    let mode = "markdown";
    let fontSize = 14;

    let editor: ace.Editor;
    let editorElement: HTMLElement;

    function bold() {
        if (editor) {
            const selectedText = editor.getSelectedText();
            const selection = editor.getSelectionRange();

            if (selectedText) {
                editor.session.replace(selection, `**${selectedText}**`);
            }
        }
    }

    function italic() {
        if (editor) {
            const selectedText = editor.getSelectedText();
            const selection = editor.getSelectionRange();

            if (selectedText) {
                editor.session.replace(selection, `*${selectedText}*`);
            }
        }
    }

    function underline() {
        if (editor) {
            const selectedText = editor.getSelectedText();
            const selection = editor.getSelectionRange();

            if (selectedText) {
                editor.session.replace(selection, `<ins>${selectedText}<ins>`);
            }
        }
    }

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

        // Sync with parent component
        editor.session.on("change", () => {
            if (appState.hurlResult) {
                appState.hurlResult.OutputString = editor.getValue();
            }
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
        if (editor && appState.hurlResult?.OutputString !== editor.getValue()) {
            editor.setValue(appState.hurlResult?.OutputString || "", -1);
        }
    });
</script>

<div bind:this={editorElement} class="flex-1 ace-editor h-full"></div>
