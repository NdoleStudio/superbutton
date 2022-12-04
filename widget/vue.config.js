const { defineConfig } = require("@vue/cli-service");
module.exports = defineConfig({
  transpileDependencies: true,
  css: {
    extract: false,
  },
  configureWebpack: {
    output: {
      filename: "widget.js",
    },
    optimization: {
      splitChunks: false,
    },
  },
  chainWebpack: (config) => {
    config.plugin("html").tap((args) => {
      args[0].inject = false;
      return args;
    });
  },
});
