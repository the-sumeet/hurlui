<script lang="ts">
	import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import { useSidebar } from "$lib/components/ui/sidebar/index.js";
	import EllipsisIcon from "@lucide/svelte/icons/ellipsis";
	import FolderIcon from "@lucide/svelte/icons/folder";
	import ShareIcon from "@lucide/svelte/icons/share";
	import Trash2Icon from "@lucide/svelte/icons/trash-2";
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
</script>

{#snippet menuItem(
	item: main.FileInfo,
	onclick: ((item: main.FileInfo) => void) | null,
)}
	<Sidebar.MenuItem>
		<Sidebar.MenuButton
			isActive={item.path === explorerState?.selectedFile?.path}
			aria-disabled={isBusy}
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
					{:else}
						<File />
					{/if}
					<span>{item.name}</span>
				</a>
			{/snippet}
		</Sidebar.MenuButton>
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
				<DropdownMenu.Separator />
				<!-- extra items removed or repurposed -->
			</DropdownMenu.Content>
		</DropdownMenu.Root>
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
