import { config } from '@vue/test-utils'
import { createPinia } from 'pinia'

// Setup Pinia for testing
config.global.plugins = [createPinia()] 