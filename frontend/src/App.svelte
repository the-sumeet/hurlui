<script lang="ts">
  import AppSidebar from "$lib/components/app-sidebar.svelte";
  import { Separator } from "$lib/components/ui/separator/index.js";
  import { Play } from "lucide-svelte";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import "./app.css";
  import * as Resizable from "$lib/components/ui/resizable/index.js";
  import Editor from "./Editor.svelte";
  import * as Select from "$lib/components/ui/select/index.js";

  import {
    GetFileContent,
    GetFiles,
    GetHurlResult,
    WriteToSelectedFile,
    GetEnvVars,
    GetEnvFilePath,
    RenamePath,
    DeletePath,
  } from "../wailsjs/go/main/App.js";
  import { main } from "../wailsjs/go/models";
  import { onMount } from "svelte";
  import Loader2Icon from "@lucide/svelte/icons/loader-2";
  import { FolderPlus } from "lucide-svelte";
  import { Save } from "lucide-svelte";
  import { Button, buttonVariants } from "$lib/components/ui/button/index.js";
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import { FilePlus } from "lucide-svelte";
  import { Info } from "lucide-svelte";
  import {
    ChangeDirectory,
    NavigateUp,
    ExecuteHurl,
    SelectFile,
    CreateNewFile,
    CreateFolder,
  } from "../wailsjs/go/main/App.js";
  import HurlReport from "./HurlReport.svelte";
  import { appState, type Dialog as AppDialog } from "./state.svelte";

  // Control the Dialog.Root via binding so closing via ESC/click-out updates state
  let dialogOpen: boolean = $state(false);
  $effect(() => {
    dialogOpen = appState.dialog != null;
  });
  $effect(() => {
    if (!dialogOpen && appState.dialog) {
      // Sync close actions from the Dialog component back to app state
      appState.dialog = null;
    }
  });

  let explorerState: main.FileExplorerState | null = $state(null);
  $effect(() => {
    if (!explorerState?.selectedFile?.path) return;

    // Todo: fetch file content only if new file is selected.
    GetFileContent(explorerState?.selectedFile?.path!).then((result) => {
      inputFileContent = result.fileContent || "";
      GetHurlResult(explorerState?.selectedFile?.path!).then((result) => {
        hurlReport = result.hurlReport || null;
      });
    });

    if (explorerState?.selectedFile) {
      GetHurlResult(explorerState?.selectedFile?.path).then((result) => {
        hurlReport = result.hurlReport || null;
      });
    }
  });
  let files: main.FileInfo[] | null = $state(null);
  let runningHurl: boolean = $state(false);
  let hurlReport: main.HurlSession[] | null = $state(null);
  // Dialog text is stored in appState.dialog.inputValue; no local mirror needed.
  let inputFileContent: string = $state("");
  let envFilePath: string = $state("");

  let envs: string[] = [];
  let selectedEnv: string = $state("");

  // Persist selected environment in localStorage
  $effect(() => {
    try {
      if (selectedEnv !== undefined) {
        localStorage.setItem("selectedEnv", selectedEnv ?? "");
      }
    } catch (e) {
      // Ignore storage errors (e.g., disabled storage)
      console.warn("Failed to persist selectedEnv:", e);
    }
  });

  function showErrorDialog(title: string, description: string) {
    appState.dialog = {
      title,
      description,
      buttonTitle: "Close",
      onclick: () => {
        appState.dialog = null;
      },
    };
  }

  async function saveSelectedFileOrDialog(): Promise<boolean> {
    if (!explorerState?.selectedFile?.path) return true;

    const res = await WriteToSelectedFile(inputFileContent);
    if (res?.error) {
      showErrorDialog("Save Error", res.error);
      return false;
    }
    return true;
  }

  function showSaveFileDialog(fileContent: string = "") {
    appState.dialog = {
      title: "Save File",
      description: `Create a new Hurl file in ${explorerState?.currentDir?.path || ""}`,
      inputLabel: "File Name",
      inputValue: "untitled.hurl",
      onclick: () => {
        const name = appState.dialog?.inputValue || "";
        console.log("Creating new file:", name);
        CreateNewFile(name, fileContent).then((result) => {
          if (result.error) {
            console.error("Failed to create new file:", result.error);
          } else {
            fetchFiles();
          }
        });
        appState.dialog = null;
      },
    };
  }

  function showNewFolderDialog() {
    appState.dialog = {
      title: "Create New Folder",
      description: `Create a new folder in ${explorerState?.currentDir?.path || ""}`,
      inputLabel: "Folder Name",
      inputValue: "NewFolder",
      onclick: () => {
        const name = appState.dialog?.inputValue || "";
        CreateFolder(name).then((result) => {
          if (result.error) {
            console.error("Failed to create new folder:", result.error);
          } else {
            fetchFiles();
          }
        });
        appState.dialog = null;
      },
    };
  }

  function showRenameDialog(item: main.FileInfo) {
    const dir = explorerState?.currentDir?.path || "";
    appState.dialog = {
      title: item.isDir ? "Rename Folder" : "Rename File",
      description: `Rename ${item.path}`,
      inputLabel: "New Name",
      inputValue: item.name,
      onclick: () => {
        const newName = appState.dialog?.inputValue || "";
        if (!newName) return;
        RenamePath(item.path, newName).then((result) => {
          if (result.error) {
            console.error("Failed to rename:", result.error);
          } else {
            // If currently selected file is renamed, update selection
            const newPath = `${dir}/${newName}`.replace(/\/+/g, "/");
            if (explorerState?.selectedFile?.path === item.path) {
              SelectFile(newPath).then((selectResult) => {
                explorerState = selectResult.fileExplorer;
              });
            }
            fetchFiles();
          }
        });
        appState.dialog = null;
      },
    };
  }

  function showDeleteDialog(item: main.FileInfo) {
    const isDir = item.isDir;
    appState.dialog = {
      title: isDir ? "Delete Folder" : "Delete File",
      description: isDir
        ? `This will permanently delete the folder and all its contents.\n${item.path}`
        : `This will permanently delete the file.\n${item.path}`,
      buttonTitle: "Delete",
      onclick: () => {
        DeletePath(item.path).then((result) => {
          if (result.error) {
            console.error("Failed to delete:", result.error);
          } else {
            fetchFiles();
          }
        });
        appState.dialog = null;
      },
    };
  }

  function onDirSelect(dir: main.FileInfo) {
    ChangeDirectory(dir.path).then(() => {
      fetchFiles();
    });
  }

  async function onFileSelect(file: main.FileInfo) {
    const saved = await saveSelectedFileOrDialog();
    if (!saved) return;
    const select_file_result = await SelectFile(file.path);
    explorerState = select_file_result.fileExplorer;
  }

  function onNavigateUp() {
    NavigateUp().then(() => {
      fetchFiles();
    });
  }

  function fetchFiles() {
    GetFiles().then((result) => {
      console.log("Fetched files:", result.files);
      explorerState = result.fileExplorer;
      files = result.files;
    });
  }

  async function onExecuteHurl() {
    if (!explorerState?.selectedFile) {
      console.error("No file selected to execute Hurl");
      return;
    }

    // First, save the current content or show error dialog.
    const saved = await saveSelectedFileOrDialog();
    if (!saved) return;

    runningHurl = true;
    ExecuteHurl(explorerState?.selectedFile?.path, selectedEnv).then(
      (result) => {
        console.log("Hurl execution result:", result);

        if (result?.error) {
          appState.dialog = {
            title: "Execution Error",
            description: result.error,
            buttonTitle: "Close",
            onclick: () => {
              appState.dialog = null;
            },
          };
          runningHurl = false;
          return;
        }

        console.log("result", result);
        hurlReport = result.hurlReport || null;
        runningHurl = false;
      },
    );
  }

  onMount(() => {
    fetchFiles();

    GetEnvVars().then((result) => {
      envs = result.envs || [];
    });

    GetEnvFilePath().then((result) => {
      envFilePath = result.envFilePath || "";
    });

    // Restore selected environment from localStorage
    try {
      const storedEnv = localStorage.getItem("selectedEnv");
      if (storedEnv !== null) {
        selectedEnv = storedEnv;
      }
    } catch (e) {
      console.warn("Failed to restore selectedEnv:", e);
    }
  });
</script>

{#snippet dialog(dialog: AppDialog)}
  <Dialog.Content class="sm:max-w-[425px]">
    <Dialog.Header>
      <Dialog.Title>{dialog.title}</Dialog.Title>
      {#if dialog.description}
        <Dialog.Description>
          {dialog.description}
        </Dialog.Description>
      {/if}
    </Dialog.Header>
    <div class="grid gap-4 py-4">
      {#if dialog.inputLabel}
        <div class="grid grid-cols-4 items-center gap-4">
          <Label for="name" class="text-right">{dialog.inputLabel}</Label>
          <Input
            id="name"
            bind:value={appState.dialog!.inputValue}
            class="col-span-3"
          />
        </div>
      {/if}
      <!-- <div class="grid grid-cols-4 items-center gap-4">
        <Label for="username" class="text-right">Username</Label>
        <Input id="username" value="@peduarte" class="col-span-3" />
      </div> -->
    </div>
    {#if dialog.onclick}
      <Dialog.Footer>
        <Button type="button" onclick={dialog.onclick} disabled={runningHurl}
          >{dialog.buttonTitle ?? "OK"}</Button
        >
      </Dialog.Footer>
    {/if}
  </Dialog.Content>
{/snippet}

<Sidebar.Provider class="h-screen overflow-y-hidden">
  <!-- Sidebar -->
  <AppSidebar
    {explorerState}
    {files}
    {onDirSelect}
    {onFileSelect}
    {onNavigateUp}
    onRename={showRenameDialog}
    onDelete={showDeleteDialog}
    isBusy={runningHurl}
    class="h-full"
  />

  <!-- Main content -->
  <Sidebar.Inset class="h-full ">
    <!-- Dialog -->
    <Dialog.Root bind:open={dialogOpen}>
      <!-- <Dialog.Trigger class={buttonVariants({ variant: "outline" })}
        >Edit Profile</Dialog.Trigger
      > -->
      {#if appState.dialog}
        {@render dialog(appState.dialog)}
      {/if}
    </Dialog.Root>

    <!-- Header -->
    <header class="flex h-16 shrink-0 items-center gap-2 border-b px-4">
      <Sidebar.Trigger class="-ml-1" disabled={runningHurl} />
      <Separator
        orientation="vertical"
        class="mr-2 data-[orientation=vertical]:h-4"
      />

      <!-- Toolbar -->
      <div class="p-1 flex w-full justify-end gap-1">
        <Button
          disabled={runningHurl ||
            !explorerState?.selectedFile.path ||
            !explorerState.selectedFile.name.endsWith(".hurl")}
          onclick={() => {
            if (runningHurl) return;

            if (!explorerState?.selectedFile.path) {
              showSaveFileDialog();
            } else {
              onExecuteHurl();
            }
          }}
          >{#if runningHurl}
            <Loader2Icon class="animate-spin" />
            Running
          {:else}
            <Play />
            Run
          {/if}</Button
        >

        <!-- <Button
          variant="outline"
          disabled={runningHurl}
          onclick={() => {
            if (!explorerState?.selectedFile?.path) {
              // No file yet — prompt to save as new
              showSaveFileDialog();
            } else {
              // Persist current buffer to selected file
              WriteToSelectedFile(inputFileContent).then((res) => {
                if (res?.error) {
                  console.error("Failed to save file:", res.error);
                }
              });
            }
          }}
        >
          <Save />
        </Button> -->
        <Button
          variant="outline"
          onclick={showNewFolderDialog}
          disabled={runningHurl}><FolderPlus /></Button
        >
        <Button
          variant="outline"
          onclick={() => showSaveFileDialog("")}
          disabled={runningHurl}
        >
          <FilePlus /></Button
        >

        <Select.Root type="single" bind:value={selectedEnv}>
          <Select.Trigger class="w-min">{selectedEnv || "Env"}</Select.Trigger>
          <Select.Content>
            <Select.Item value="">None</Select.Item>
            {#each envs as env}
              <Select.Item value={env}>{env}</Select.Item>
            {/each}
            <Select.Item disabled value="NOT_SET"
              >Env vars in {envFilePath}</Select.Item
            >
          </Select.Content>
        </Select.Root>
      </div>

      <!-- Breadcrumbs -->
      <!-- <Breadcrumb.Root>
        <Breadcrumb.List>
          <Breadcrumb.Item class="hidden md:block">
            <Breadcrumb.Link href="#">lib</Breadcrumb.Link>
          </Breadcrumb.Item>
          <Breadcrumb.Separator class="hidden md:block" />
          <Breadcrumb.Item class="hidden md:block">
            <Breadcrumb.Link href="#">components</Breadcrumb.Link>
          </Breadcrumb.Item>
          <Breadcrumb.Separator class="hidden md:block" />
          <Breadcrumb.Item>
            <Breadcrumb.Page>button.svelte</Breadcrumb.Page>
          </Breadcrumb.Item>
        </Breadcrumb.List>
      </Breadcrumb.Root> -->
    </header>

    <Resizable.PaneGroup
      direction="horizontal"
      class="min-h-[200px] h-full border overflow-y-hidden flex-1"
    >
      <!-- Input -->
      <Resizable.Pane defaultSize={50} class="h-full">
        <Editor bind:content={inputFileContent} />
      </Resizable.Pane>

      <!-- Output -->
      {#if hurlReport && hurlReport.length > 0}
        <Resizable.Handle withHandle />
        <Resizable.Pane defaultSize={50} class="h-full overflow-y-hidden">
          <HurlReport {hurlReport} />
        </Resizable.Pane>
      {/if}
    </Resizable.PaneGroup>
  </Sidebar.Inset>
</Sidebar.Provider>
