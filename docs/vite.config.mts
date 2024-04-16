import { join } from 'node:path'
import { defineConfig } from 'vite'
import UnoCSS from 'unocss/vite'
import Inspect from 'vite-plugin-inspect'

import { GitChangelog,GitChangelogMarkdownSection } from '@nolebase/vitepress-plugin-git-changelog/vite'

export default defineConfig({
  ssr: {
    noExternal: [
      // If there are other packages that need to be processed by Vite, you can add them here.
      '@nolebase/vitepress-plugin-enhanced-readabilities',
      '@nolebase/vitepress-plugin-highlight-targeted-heading',
      'asciinema-player'
    ],
  },
  optimizeDeps: {
    exclude: ['vitepress'],
  },
  plugins: [
    Inspect(),
    // https://vitejs.dev/guide/api-plugin.html
    GitChangelog({
      repoURL: 'https://github.com/nekomeowww/ollama-operator',
      rewritePaths: {
        'docs/': ''
      }
    }),
    GitChangelogMarkdownSection({
      locales: {
        'en': {
          gitChangelogMarkdownSectionTitles: {
            changelog: 'Changelog',
            contributors: 'Contributors'
          }
        },
        'zh-CN': {
          gitChangelogMarkdownSectionTitles: {
            changelog: '页面历史',
            contributors: '贡献者'
          }
        }
      },
      excludes: [
        join('pages', 'en', 'index.md'),
        join('pages', 'zh-CN', 'index.md')
      ],
    }),
    UnoCSS()
  ],
})
