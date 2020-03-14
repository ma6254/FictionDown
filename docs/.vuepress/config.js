module.exports = {
    title: 'FictionDown',
    description: 'FictionDown 小说爬取',
    markdown: {
        lineNumbers: true
    },
    themeConfig: {
        repo: 'ma6254/FictionDown',
        docsDir: 'docs',
        editLinks: true,
        editLinkText: '帮助我们改善此页面！',
        repoLabel: '查看源码',
        lastUpdated: '上次更新',
        smoothScroll: true,
        sidebar: [
            {
                title: '指南',
                path: '/guide/',
                collapsable: false,
                sidebarDepth: 3,
                children: [
                    '/guide/install',
                    '/guide/quickstart',
                    '/guide/source'
                ]
            },
            {
                title: 'API',
                path: '/api/',
                collapsable: false,
                sidebarDepth: 3,
            },
        ],
        nav: [
            { text: '指南', link: '/guide/' },
            { text: 'API', link: '/api/' },
        ]
    }
}
