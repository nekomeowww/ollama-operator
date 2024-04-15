// https://vitepress.dev/guide/custom-theme
import { h } from 'vue'
import type { Theme } from 'vitepress'
import DefaultTheme from 'vitepress/theme'

import { InjectionKey, NolebaseEnhancedReadabilitiesMenu, NolebaseEnhancedReadabilitiesScreenMenu } from '@nolebase/vitepress-plugin-enhanced-readabilities/client'
import { NolebaseGitChangelogPlugin } from '@nolebase/vitepress-plugin-git-changelog/client'
import {
  NolebaseHighlightTargetedHeading,
} from '@nolebase/vitepress-plugin-highlight-targeted-heading/client'


import '@nolebase/vitepress-plugin-enhanced-mark/client/style.css'
import '@nolebase/vitepress-plugin-enhanced-readabilities/client/style.css'
import '@nolebase/vitepress-plugin-highlight-targeted-heading/client/style.css'

import AsciinemaPlayer from './components/AsciinemaPlayer.vue'

import './style.css'

export default {
  extends: DefaultTheme,
  Layout: () => {
    return h(DefaultTheme.Layout, null, {
      'layout-top': () => [
        h(NolebaseHighlightTargetedHeading),
      ],
      // A enhanced readabilities menu for wider screens
      'nav-bar-content-after': () => [
        h(NolebaseEnhancedReadabilitiesMenu)
      ],
      // A enhanced readabilities menu for narrower screens (usually smaller than iPad Mini)
      'nav-screen-content-after': () => [
        h(NolebaseEnhancedReadabilitiesScreenMenu)
      ],
    })
  },
  enhanceApp({ app }) {
    app.component('AsciinemaPlayer', AsciinemaPlayer)
    app.provide(InjectionKey, {
      spotlight: {
        defaultToggle: true
      }
    })
    app.use(NolebaseGitChangelogPlugin)
  }
} satisfies Theme
