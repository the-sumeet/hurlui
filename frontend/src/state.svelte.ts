import { main } from "../wailsjs/go/models"


export interface Dialog {
    title: string
    description?: string | null
    buttonTitle?: string | null
    inputLabel?: string | null
    inputValue?: string | null
    onclick?: () => void | null
    validator?: (value: string) => string | null
}
interface AppState {
    hurlResult: main.HurlResult | null
    dialog: Dialog | null
}

export const appState: AppState = $state({
    hurlResult: null,
    dialog: null
})

