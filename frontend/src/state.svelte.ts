import { main } from "../wailsjs/go/models"

interface AppState {
    hurlResult: main.HurlResult | null
}

export const appState: AppState = $state({
    hurlResult: null
})