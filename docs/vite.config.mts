import { join } from 'node:path'
import { defineConfig } from 'vite'
import UnoCSS from 'unocss/vite'
import Inspect from 'vite-plugin-inspect'

import { GitChangelog, GitChangelogMarkdownSection } from '@nolebase/vitepress-plugin-git-changelog/vite'

export default defineConfig({
  ssr: {
    noExternal: [
      // If there are other packages that need to be processed by Vite, you can add them here.
      '@nolebase/vitepress-plugin-enhanced-readabilities',
      '@nolebase/vitepress-plugin-highlight-targeted-heading',
      '@nolebase/ui-asciinema',
    ],
  },
  optimizeDeps: {
    exclude: [
      'vitepress',
    ],
  },
  plugins: [
    Inspect(),
    // https://vitejs.dev/guide/api-plugin.html
    GitChangelog({
      repoURL: 'https://github.com/nekomeowww/ollama-operator',
    }),
    GitChangelogMarkdownSection({
      excludes: [
        join('pages', 'en', 'index.md'),
        join('pages', 'zh-CN', 'index.md'),
      ],
    }),
    UnoCSS(),
  ],
})
