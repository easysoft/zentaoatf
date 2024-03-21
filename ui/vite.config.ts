import { resolve } from 'path'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'


// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      '@': resolve('.', './src'),
      '~': resolve('.', './node_modules'),
    },
    extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json', '.less'],
  },
  build: {
    outDir: '../client/ui',
  },
  plugins: [vue()],
})
