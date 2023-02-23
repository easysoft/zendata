const path = require('path');

function resolve(dir) {
    return path.join(__dirname, dir)
}

module.exports = {
    // productionSourceMap: false,
    publicPath: process.env.NODE_ENV === 'production' && !process.env.UI_IN_CLIENT ? 'ui': '',

    chainWebpack: (config) => {
        config.resolve.alias
            .set('@', resolve('./src'))
        // .set('components',resolve('./src/components'))
    }
}
