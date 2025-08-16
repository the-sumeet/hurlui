<script lang="ts">
  import AppSidebar from "$lib/components/app-sidebar.svelte";
  import * as Breadcrumb from "$lib/components/ui/breadcrumb/index.js";
  import { Separator } from "$lib/components/ui/separator/index.js";
  import { Play } from "lucide-svelte";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import "./app.css";
  import { Button } from "$lib/components/ui/button/index.js";
  import * as Resizable from "$lib/Components/ui/resizable/index.js";
  import Editor from "./Editor.svelte";
  import Output from "./Output.svelte";
  let resultText: string = "Please enter your name below 👇";
  let name: string;
  import { GetFiles, GetHurlResult } from "../wailsjs/go/main/App.js";
  import { main } from "../wailsjs/go/models";
  import { onMount } from "svelte";
  import Loader2Icon from "@lucide/svelte/icons/loader-2";
  import {
    ChangeDirectory,
    NavigateUp,
    ExecuteHurl,
    SelectFile,
  } from "../wailsjs/go/main/App.js";
  import { appState } from "./state.svelte";
  import HurlReport from "./HurlReport.svelte";

  let explorerState: main.FileExplorerState | null = $state(null);
  let files: main.FileInfo[] | null = $state(null);
  let runningHurl: boolean = $state(false);
  let hurlReport: main.HurlSession[] | null = $state(null);

  function onDirSelect(dir: main.FileInfo) {
    ChangeDirectory(dir.path).then(() => {
      fetchFiles();
    });
  }

  function onFileSelect(file: main.FileInfo) {
    SelectFile(file.path).then((result) => {
      explorerState = result.fileExplorer;
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
  });
</script>

<Sidebar.Provider class="h-screen">
  <!-- Sidebar -->
  <AppSidebar
    {explorerState}
    {files}
    {onDirSelect}
    {onFileSelect}
    {onNavigateUp}
    {onExecuteHurl}
  />

  <!-- Main content -->
  <Sidebar.Inset class="h-full">
    <header class="flex h-16 shrink-0 items-center gap-2 border-b px-4">
      <Sidebar.Trigger class="-ml-1" />
      <Separator
        orientation="vertical"
        class="mr-2 data-[orientation=vertical]:h-4"
      />
      <Breadcrumb.Root>
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
      </Breadcrumb.Root>
    </header>

    <div class="flex flex-1 flex-col h-full w-full">
      <!-- Topbar -->
      <div class="p-1 flex">
        <Button onclick={onExecuteHurl}
          >{#if runningHurl}
            <Loader2Icon class="animate-spin" />
            Running
          {:else}
            <Play /> Run
          {/if}</Button
        >
      </div>

      <Resizable.PaneGroup
        direction="horizontal"
        class="min-h-[200px] h-full border overflow-y-hidden flex-1"
      >
        <!-- Input -->
        <Resizable.Pane defaultSize={50} class="h-full">
          <Editor />
        </Resizable.Pane>

        <!-- Output -->
        {#if hurlReport && hurlReport.length > 0}
          <Resizable.Handle withHandle />
          <Resizable.Pane defaultSize={50} class="h-full overflow-y-hidden">
            <HurlReport {hurlReport} />
          </Resizable.Pane>
        {/if}
      </Resizable.PaneGroup>
    </div>
  </Sidebar.Inset>
</Sidebar.Provider>
