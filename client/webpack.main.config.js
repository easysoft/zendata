// const CopyPlugin = require("copy-webpack-plugin");

module.exports = {
  entry: './src/main.js',
  module: {
    rules: require('./webpack.rules'),
  },
  // plugins: [
  //     new CopyPlugin({
  //       patterns: [{ from: "./icon", to: "icon" }]
  //     })
  // ],
};
