import { config } from '@vue/test-utils'
import { createPinia } from 'pinia'
import * as Vue from 'vue'

// Make Vue available globally
global.Vue = Vue

// Create a fresh Pinia instance for each test
const pinia = createPinia()
config.global.plugins = [pinia]

// Mock window.alert
global.alert = jest.fn()

// Mock console methods
console.error = jest.fn()
console.log = jest.fn()
console.warn = jest.fn() 