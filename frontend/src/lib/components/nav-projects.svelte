<script lang="ts">
	import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import { useSidebar } from "$lib/components/ui/sidebar/index.js";
	import EllipsisIcon from "@lucide/svelte/icons/ellipsis";
	import { ArrowDown } from "lucide-svelte";
	import { ArrowDownUp } from "lucide-svelte";
	import { main } from "../../../wailsjs/go/models";
	import { Folder, File } from "lucide-svelte";
	let {
		explorerState,
		files,
		onDirSelect,
		onFileSelect,
		onNavigateUp,
		onRename,
		onDelete,
		isBusy = false,
	}: {
		explorerState?: main.FileExplorerState | null;
		files?: main.FileInfo[] | null;
		onDirSelect: (dir: main.FileInfo) => void;
		onFileSelect: (file: main.FileInfo) => void;
		onNavigateUp: () => void;
		onRename: (item: main.FileInfo) => void;
		onDelete: (item: main.FileInfo) => void;
		isBusy?: boolean;
	} = $props();

	$inspect(files);

	const sidebar = useSidebar();

	function isMarkdownFile(file: main.FileInfo): boolean {
		return file.name.endsWith(".md") || file.name.endsWith(".markdown");
	}
	function isHurlFile(file: main.FileInfo): boolean {
		return file.name.endsWith(".hurl");
	}

	function canOpenFile(file: main.FileInfo): boolean {
		return file.isDir || isMarkdownFile(file) || isHurlFile(file);
	}
</script>

{#snippet menuItem(
	item: main.FileInfo,
	onclick: ((item: main.FileInfo) => void) | null,
)}
	<Sidebar.MenuItem>
		<Sidebar.MenuButton
			isActive={item.path === explorerState?.selectedFile?.path &&
				canOpenFile(item)}
			aria-disabled={isBusy || !canOpenFile(item)}
			onclick={() => {
				if (isBusy) return;
				if (onclick) {
					onclick(item);
				} else if (item.isDir) {
					onDirSelect(item);
				} else {
					onFileSelect(item);
				}
			}}
		>
			{#snippet child({ props })}
				<a href="#!" {...props} title={item.path}>
					{#if item.isDir}
						<Folder />
					{:else if isHurlFile(item)}
						<ArrowDownUp color={"#f80288"} />
					{:else if isMarkdownFile(item)}
						<ArrowDown color={"#34a7ff"} />
					{:else}
						<File />
					{/if}
					<span>{item.name}</span>
				</a>
			{/snippet}
		</Sidebar.MenuButton>
		{#if canOpenFile(item)}
			<DropdownMenu.Root>
				<DropdownMenu.Trigger>
					{#snippet child({ props })}
						<Sidebar.MenuAction
							showOnHover
							{...props}
							aria-disabled={isBusy}
						>
							<EllipsisIcon />
							<span class="sr-only">More</span>
						</Sidebar.MenuAction>
					{/snippet}
				</DropdownMenu.Trigger>
				<DropdownMenu.Content
					class="w-48"
					side={sidebar.isMobile ? "bottom" : "right"}
					align={sidebar.isMobile ? "end" : "start"}
				>
					<DropdownMenu.Item onclick={() => onRename(item)}>
						<span>Rename</span>
					</DropdownMenu.Item>
					<DropdownMenu.Item onclick={() => onDelete(item)}>
						<span>Delete</span>
					</DropdownMenu.Item>
					<!-- <DropdownMenu.Separator /> -->
					<!-- extra items removed or repurposed -->
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		{/if}
	</Sidebar.MenuItem>
{/snippet}

<Sidebar.Group class="group-data-[collapsible=icon]:hidden">
	<!-- <Sidebar.GroupLabel>Projects</Sidebar.GroupLabel> -->
	<Sidebar.Menu>
		<!-- Back button -->
		{#if explorerState?.currentDir.path !== "/"}
			{@render menuItem(
				{
					name: "..",
					isDir: true,
					path: "NA",
					size: 0,
					modified: new Date().toISOString(),
				},
				() => {
					onNavigateUp();
				},
			)}
		{/if}

		{#each files || [] as item (item.path)}
			{@render menuItem(item, null)}
		{/each}
	</Sidebar.Menu>
</Sidebar.Group>
