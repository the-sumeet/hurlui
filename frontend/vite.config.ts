import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import path from "path";
// https://vitejs.dev/config/
export default defineConfig({
    resolve: {
        alias: {
            $lib: path.resolve("./src/lib"),
        },
    }, plugins: [tailwindcss(), svelte()]
});
