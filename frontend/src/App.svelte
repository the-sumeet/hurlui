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

  let dialogOpened = $derived(appState.dialog != null);

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
  let dialogInput: string = $state("");
  let inputFileContent: string = $state("");
  $effect(() => {
    if (explorerState?.selectedFile) {
      WriteToSelectedFile(inputFileContent);
    }
  });

  let envs: string[] = [];
  let selectedEnv: string = $state("");

  function showSaveFileDialog(fileContent: string = "") {
    if (fileContent == "") {
      fileContent = inputFileContent;
    }

    dialogInput = "untitled.hurl";
    appState.dialog = {
      title: "Save File",
      description: `Create a new Hurl file in ${explorerState?.currentDir?.path || ""}`,
      inputLabel: "File Name",
      onclick: () => {
        if (!dialogInput.endsWith(".hurl")) {
          dialogInput += ".hurl";
        }
        // Handle file save logic here
        console.log("Creating new file:", dialogInput);
        CreateNewFile(dialogInput, fileContent).then((result) => {
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
    dialogInput = "NewFolder";
    appState.dialog = {
      title: "Create New Folder",
      description: `Create a new folder in ${explorerState?.currentDir?.path || ""}`,
      inputLabel: "Folder Name",
      onclick: () => {
        CreateFolder(dialogInput).then((result) => {
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

  function onDirSelect(dir: main.FileInfo) {
    ChangeDirectory(dir.path).then(() => {
      fetchFiles();
    });
  }

  function onFileSelect(file: main.FileInfo) {
    SelectFile(file.path).then((select_file_result) => {
      explorerState = select_file_result.fileExplorer;

      // GetFileContent(file.path).then((result) => {
      //   inputFileContent = result.fileContent || "";
      // });
    });
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

  function onExecuteHurl() {
    if (!explorerState?.selectedFile) {
      console.error("No file selected to execute Hurl");
      return;
    }
    runningHurl = true;
    ExecuteHurl(explorerState?.selectedFile?.path).then((result) => {
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

      if (explorerState?.selectedFile) {
        GetHurlResult(explorerState?.selectedFile?.path).then((result) => {
          hurlReport = result.hurlReport || null;
        });
      }
      runningHurl = false;
    });
  }

  onMount(() => {
    fetchFiles();

    GetEnvVars().then((result) => {
      envs = result.envs || [];
    });
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
          <Input id="name" bind:value={dialogInput} class="col-span-3" />
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
    isBusy={runningHurl}
    class="h-full"
  />

  <!-- Main content -->
  <Sidebar.Inset class="h-full ">
    <!-- Dialog -->
    <Dialog.Root open={dialogOpened}>
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

        <Button
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
        </Button>
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
