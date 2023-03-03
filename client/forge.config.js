module.exports = {
    electronPackagerConfig: {
        "name": "zd1",
        "icon": "./ui/favicon.ico"
    },
    packagerConfig: {
        "name": "zd1",
        "icon": "./icon/favicon.icns",
        extraResource: [
            './bin',
            './ui',
            './lang',
        ]
    },
    makers: [
        {
            name: '@electron-forge/maker-squirrel',
            config: {
                name: 'zd'
            }
        },
        {
            name: '@electron-forge/maker-zip',
            platforms: [
                'darwin'
            ]
        },
        {
            name: '@electron-forge/maker-deb',
            config: {}
        },
        {
            name: '@electron-forge/maker-rpm',
            config: {}
        }
    ],
    plugins: [
        [
            "@electron-forge/plugin-webpack",
            {
                mainConfig: './webpack.main.config.js',
                renderer: {
                    config: './webpack.renderer.config.js',
                    entryPoints: [
                        // {
                        //   html: './src/index.html',
                        //   js: './src/renderer.js',
                        //   name: 'main_window'
                        // }
                    ]
                }
            }
        ],
        [
            "@timfish/forge-externals-plugin",
            {
                "externals": ["@electron/remote"],
                "includeDeps": true
            }
        ]
    ]
}
