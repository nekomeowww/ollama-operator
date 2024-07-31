// https://vitepress.dev/guide/custom-theme
import { h } from 'vue'
import type { Theme } from 'vitepress'
import DefaultTheme from 'vitepress/theme'

import {
  InjectionKey,
  NolebaseEnhancedReadabilitiesMenu,
  NolebaseEnhancedReadabilitiesScreenMenu,
} from '@nolebase/vitepress-plugin-enhanced-readabilities/client'
import {
  NolebaseGitChangelogPlugin,
} from '@nolebase/vitepress-plugin-git-changelog/client'
import {
  NolebaseHighlightTargetedHeading,
} from '@nolebase/vitepress-plugin-highlight-targeted-heading/client'

import { enhanceAppWithTabs } from 'vitepress-plugin-tabs/client'

import TitleBlockContainer from './components/TitleBlockContainer.vue'
import TitleBlockContainerGroup from './components/TitleBlockContainerGroup.vue'
import GettingStartedBlocksEn from './components/GettingStartedBlocksEn.vue'
import GettingStartedBlocksZhCn from './components/GettingStartedBlocksZhCn.vue'

import '@nolebase/vitepress-plugin-enhanced-mark/client/style.css'
import '@nolebase/vitepress-plugin-enhanced-readabilities/client/style.css'
import '@nolebase/vitepress-plugin-highlight-targeted-heading/client/style.css'

import 'asciinema-player/dist/bundle/asciinema-player.css'

import 'virtual:uno.css'
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
        h(NolebaseEnhancedReadabilitiesMenu),
      ],
      // A enhanced readabilities menu for narrower screens (usually smaller than iPad Mini)
      'nav-screen-content-after': () => [
        h(NolebaseEnhancedReadabilitiesScreenMenu),
      ],
    })
  },
  enhanceApp({ app }) {
    app.component('TitleBlockContainer', TitleBlockContainer)
    app.component('TitleBlockContainerGroup', TitleBlockContainerGroup)
    app.component('GettingStartedBlocksEn', GettingStartedBlocksEn)
    app.component('GettingStartedBlocksZhCn', GettingStartedBlocksZhCn)

    app.use(enhanceAppWithTabs)
    app.provide(InjectionKey, {
      spotlight: {
        defaultToggle: true,
      },
    })
    app.use(NolebaseGitChangelogPlugin)
  },
} satisfies Theme
