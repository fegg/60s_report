import { defineConfig } from 'vite';
import tsconfigPaths from 'vite-tsconfig-paths';
import { viteMockServe } from "vite-plugin-mock";

export default defineConfig({
  plugins: [
    tsconfigPaths(),
    viteMockServe()
  ],
  build: {
    outDir: "build"
  },
  server: {
    host: '0.0.0.0'
  }
})