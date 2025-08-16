import './app.css'
import App from './App.svelte'
import { mount } from 'svelte';

// Import bootstrap icons
const app = mount(App, { target: document.getElementById("app") });

export default app

