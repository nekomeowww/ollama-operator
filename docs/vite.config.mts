import { join } from 'node:path'
import { defineConfig } from 'vite'
import UnoCSS from 'unocss/vite'

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
    // https://vitejs.dev/guide/api-plugin.html
    GitChangelog({
      repoURL: 'https://github.com/nekomeowww/ollama-operator',
      rewritePaths: {
        'docs/': ''
      }
    }),
    GitChangelogMarkdownSection({
      getChangelogTitle: (_, __, { helpers }): string => {
        if (helpers.idStartsWith(join('pages', 'en')))
          return 'Page History'
        if (helpers.idStartsWith(join('pages', 'zh-CN')))
          return '页面历史'

        return 'Page History'
      },
      getContributorsTitle: (_, __, { helpers }): string => {
        if (helpers.idStartsWith(join('pages', 'en')))
          return 'Contributors'
        if (helpers.idStartsWith(join('pages', 'zh-CN')))
          return '贡献者'

        return 'Contributors'
      },
      excludes: [],
      exclude: (_, { helpers }): boolean => {
        if (helpers.idEquals(join('pages', 'en', 'index.md')))
          return true
        if (helpers.idEquals(join('pages', 'zh-CN', 'index.md')))
          return true

        return false
      },
    }),
    UnoCSS()
  ],
})
