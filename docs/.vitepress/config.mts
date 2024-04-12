import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Ollama Operator",
  description: "Yet another operator for running large language models on Kubernetes with ease. 🙀",
  lastUpdated: true,
  ignoreDeadLinks: [
    // Site Config | VitePress
    // https://vitepress.dev/reference/site-config#ignoredeadlinks
    //
    // ignore all localhost links
    /^https?:\/\/localhost/,
  ],
  themeConfig: {
    outline: 'deep',
    socialLinks: [
      { icon: 'github', link: 'https://github.com/nekomeowww/ollama-operator' }
    ],
    search: {
      provider: 'local',
      options: {
        locales: {
          'zh-CN': {
            translations: {
              button: {
                buttonText: '搜索文档',
                buttonAriaLabel: '搜索文档',
              },
              modal: {
                noResultsText: '无法找到相关结果',
                resetButtonTitle: '清除查询条件',
                footer: {
                  selectText: '选择',
                  navigateText: '切换',
                },
              },
            },
          },
        },
      },
    },
  },
  head: [
    ['link', { rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-touch-icon.png'}],
    ['link', { rel: 'icon',type: 'image/png', href:'/logo.png'}],
    ['link', { rel: 'manifest', href: '/site.webmanifest'}],
    ['link', { rel: 'mask-icon', href: '/safari-pinned-tab.svg', color:'#5bbad5'}],
    ['meta', { name: 'msapplication-TileColor', content: '#2b5797'}],
    ['meta', { name: 'theme-color', content: '#ffffff'}],
  ],
  locales: {
    'root': {
      label: 'English',
      lang: 'en',
      link: '/pages/en/',
      title: 'Ollama Operator',
      description: 'Yet another operator for running large language models on Kubernetes with ease. 🙀',
      themeConfig: {
        nav: [
          {
            text: 'Guide',
            items: [
              { text: 'Overview', link: '/pages/en/guide/overview' },
              { text: 'Supported models', link: '/pages/en/guide/supported-models' },
              { text: 'Getting started', link: '/pages/en/guide/getting-started' },
            ]
          },
          {
            text: 'Reference',
            items: [
              { text: 'CRD definition', link: '/pages/en/references/crd' },
              { text: 'Architectural Design', link: '/pages/en/references/architectural-design' },
            ]
          }
        ],
        sidebar: [
          {
            text: 'Guide',
            items: [
              { text: 'Overview', link: '/pages/en/guide/overview' },
              { text: 'Supported models', link: '/pages/en/guide/supported-models' },
              { text: 'Getting started', link: '/pages/en/guide/getting-started' },
            ]
          },
          {
            text: 'Reference',
            items: [
              { text: 'CRD definition', link: '/pages/en/references/crd' },
              { text: 'Architectural Design', link: '/pages/en/references/architectural-design' },
            ]
          }
        ]
      },
    },
    'zh-CN': {
      label: '简体中文',
      lang: 'zh-CN',
      link: '/pages/zh-CN/',
      title: 'Ollama Operator',
      description: '一个用于在 Kubernetes 上轻松运行大型语言模型的 Operator。 🙀',
      themeConfig: {
        nav: [
          {
            text: '指南',
            items: [
              { text: '概览', link: '/pages/zh-CN/guide/overview' },
              { text: '支持模型', link: '/pages/zh-CN/guide/supported-models' },
              { text: '快速开始', link: '/pages/zh-CN/guide/getting-started' },
            ]
          },
          {
            text: '参考',
            items: [
              { text: 'CRD 定义', link: '/pages/zh-CN/references/crd' },
              { text: '架构设计', link: '/pages/zh-CN/references/architectural-design' },
            ]
          }
        ],
        sidebar: [
          {
            text: '指南',
            items: [
              { text: '概览', link: '/pages/zh-CN/guide/overview' },
              { text: '支持模型', link: '/pages/zh-CN/guide/supported-models' },
              { text: '快速开始', link: '/pages/zh-CN/guide/getting-started' },
            ]
          },
          {
            text: '参考',
            items: [
              { text: 'CRD 定义', link: '/pages/zh-CN/references/crd' },
              { text: '架构设计', link: '/pages/zh-CN/references/architectural-design' },
            ]
          }
        ]
      },
    }
  }
})
