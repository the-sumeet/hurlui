<script lang="ts">
	import NavProjects from "./nav-projects.svelte";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import CommandIcon from "@lucide/svelte/icons/command";
	import type { ComponentProps } from "svelte";
	import type { main } from "wailsjs/go/models";

	interface Props extends ComponentProps<typeof Sidebar.Root> {
		ref?: any;
		explorerState?: main.FileExplorerState | null;
		files?: main.FileInfo[] | null;
		onDirSelect: (dir: main.FileInfo) => void;
		onFileSelect: (file: main.FileInfo) => void;
		onNavigateUp: () => void;
		onRename: (item: main.FileInfo) => void;
		onDelete: (item: main.FileInfo) => void;
		isBusy?: boolean;
		[key: string]: any;
	}

	let {
		ref = $bindable(null),
		explorerState,
		files,
		onDirSelect,
		onFileSelect,
		onNavigateUp,
		onRename,
		onDelete,
		isBusy = false,
		...restProps
	}: Props = $props();
</script>

<Sidebar.Root bind:ref variant="inset" {...restProps}>
	<!-- <Sidebar.Header>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton size="lg">
					{#snippet child({ props })}
						<a href="##" {...props}>
							<div
								class="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg"
							>
								<CommandIcon class="size-4" />
							</div>
							<div
								class="grid flex-1 text-left text-sm leading-tight"
							>
								<span class="truncate font-medium"
									>Acme Inc</span
								>
								<span class="truncate text-xs">Enterprise</span>
							</div>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Header> -->
	<Sidebar.Content>
		<!-- <NavMain items={data.navMain} /> -->
		<NavProjects
			{explorerState}
			{files}
			{onDirSelect}
			{onFileSelect}
			{onNavigateUp}
			onRename={onRename}
			onDelete={onDelete}
			isBusy={isBusy}
		/>
		<!-- <NavSecondary items={data.navSecondary} class="mt-auto" /> -->
	</Sidebar.Content>
	<!-- <Sidebar.Footer>
		<NavUser user={data.user} />
	</Sidebar.Footer> -->
</Sidebar.Root>
