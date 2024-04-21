import { text } from 'stream/consumers'
import { defineConfig } from 'vitepress'
import { tabsMarkdownPlugin } from 'vitepress-plugin-tabs'

// https://vitepress.dev/references/site-config
export default defineConfig({
  markdown: {
    config(md) {
      md.use(tabsMarkdownPlugin)
    }
  },
  title: "Ollama Operator",
  description: "Large language models, scaled, deployed. - Yet another operator for running large language models on Kubernetes with ease. 🙀",
  lastUpdated: true,
  ignoreDeadLinks: [
    // Site Config | VitePress
    // https://vitepress.dev/references/site-config#ignoredeadlinks
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
      description: 'Large language models, scaled, deployed - Yet another operator for running large language models on Kubernetes with ease. 🙀',
      themeConfig: {
        nav: [
          {
            text: 'Guide',
            items: [
              { text: 'Overview', link: '/pages/en/guide/overview' },
              {
                text: 'Getting started',
                items: [
                  { text: 'Install Ollama Operator', link: '/pages/en/guide/getting-started/' },
                  { text: 'Deploy models through CLI', link: '/pages/en/guide/getting-started/cli' },
                  { text: 'Deploy models through CRD', link: '/pages/en/guide/getting-started/crd' },
                ]
              },
              { text: 'Supported models', link: '/pages/en/guide/supported-models' },
            ]
          },
          {
            text: 'Reference',
            items: [
              { text: 'CLI Reference', link: '/pages/en/references/cli/' },
              { text: 'CRD Reference', link: '/pages/en/references/crd/' },
              { text: 'Architectural Design', link: '/pages/en/references/architectural-design' },
            ]
          },
          {
            text: 'Acknowledgements',
            link: '/pages/en/acknowledgements'
          }
        ],
        sidebar: [
          {
            text: 'Guide',
            items: [
              { text: 'Overview', link: '/pages/en/guide/overview' },
              {
                text: 'Getting started',
                items: [
                  { text: 'Install Ollama Operator', link: '/pages/en/guide/getting-started/' },
                  { text: 'Deploy models through CLI', link: '/pages/en/guide/getting-started/cli' },
                  { text: 'Deploy models through CRD', link: '/pages/en/guide/getting-started/crd' },
                ]
              },
              { text: 'Supported models', link: '/pages/en/guide/supported-models' },
            ]
          },
          {
            text: 'Reference',
            items: [
              {
                text: 'CLI Reference',
                items: [
                  { text: 'Commands list', link: '/pages/en/references/cli/' },
                  {
                    text: 'Commands',
                    items: [
                      { text: 'kollama deploy', link: '/pages/en/references/cli/commands/deploy' },
                      { text: 'kollama undeploy', link: '/pages/en/references/cli/commands/undeploy' },
                      { text: 'kollama expose', link: '/pages/en/references/cli/commands/undeploy' }
                    ]
                  }
                ]
              },
              {
                text: 'CRD Reference',
                items: [
                  { text: 'CRD list', link: '/pages/en/references/crd/' },
                  { text: 'Model', link: '/pages/en/references/crd/model' }
                ]
              },
              { text: 'Architectural Design', link: '/pages/en/references/architectural-design' },
            ]
          },
          {
            text: 'Acknowledgements',
            link: '/pages/en/acknowledgements'
          }
        ]
      },
    },
    'zh-CN': {
      label: '简体中文',
      lang: 'zh-CN',
      link: '/pages/zh-CN/',
      title: 'Ollama Operator',
      description: '大语言模型，伸缩自如，轻松部署 - 一个用于在 Kubernetes 上轻松运行大型语言模型的 Operator。 🙀',
      themeConfig: {
        nav: [
          {
            text: '指南',
            items: [
              { text: '概览', link: '/pages/zh-CN/guide/overview' },
              {
                text: '快速上手',
                items: [
                  { text: '安装 Ollama Operator', link: '/pages/zh-CN/guide/getting-started/' },
                  { text: '通过 CLI 部署模型', link: '/pages/zh-CN/guide/getting-started/cli' },
                  { text: '通过 CRD 部署模型', link: '/pages/zh-CN/guide/getting-started/crd' },
                ]
              },
              { text: '支持模型', link: '/pages/zh-CN/guide/supported-models' },
            ]
          },
          {
            text: '参考',
            items: [
              {
                text: 'CLI 参考',
                items: [
                  { text: '命令列表', link: '/pages/zh-CN/references/cli/' },
                ]
              },
              { text: 'CRD 参考', link: '/pages/zh-CN/references/crd/' },
              { text: '架构设计', link: '/pages/zh-CN/references/architectural-design' },
            ]
          },
          {
            text: '致谢',
            link: '/pages/zh-CN/acknowledgements'
          }
        ],
        sidebar: [
          {
            text: '指南',
            items: [
              { text: '概览', link: '/pages/zh-CN/guide/overview' },
              {
                text: '快速上手',
                items: [
                  { text: '安装 Ollama Operator', link: '/pages/zh-CN/guide/getting-started/' },
                  { text: '通过 CLI 部署模型', link: '/pages/zh-CN/guide/getting-started/cli' },
                  { text: '通过 CRD 部署模型', link: '/pages/zh-CN/guide/getting-started/crd' },
                ]
              },
              { text: '支持模型', link: '/pages/zh-CN/guide/supported-models' },
            ]
          },
          {
            text: '参考',
            items: [
              {
                text: 'CLI 参考',
                items: [
                  { text: '命令列表', link: '/pages/zh-CN/references/cli/' },
                  {
                    text: '子命令',
                    items: [
                      { text: 'kollama deploy', link: '/pages/zh-CN/references/cli/commands/deploy' },
                      { text: 'kollama undeploy', link: '/pages/zh-CN/references/cli/commands/undeploy' },
                      { text: 'kollama expose', link: '/pages/zh-CN/references/cli/commands/undeploy' }
                    ]
                  }
                ]
              },
              {
                text: 'CRD 参考',
                items: [
                  { text: 'CRD 列表', link: '/pages/zh-CN/references/crd/' },
                  { text: 'Model 模型资源', link: '/pages/zh-CN/references/crd/model' }
                ]
              },
              { text: '架构设计', link: '/pages/zh-CN/references/architectural-design' },
            ]
          },
          {
            text: '致谢',
            link: '/pages/zh-CN/acknowledgements'
          }
        ]
      },
    }
  }
})
