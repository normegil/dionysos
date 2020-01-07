/* eslint-disable */
const path = require('path')
const VueLoaderPlugin = require('vue-loader/lib/plugin')

function getConfiguration (env) {
  let cfg = {
    entry: './src/index.ts',
    mode: 'development',
    optimization: {
      usedExports: true,
    },
    output: {
      filename: 'bundle.js',
      path: path.resolve(__dirname, 'dist'),
    },
    module: {
      rules: [
        {
          test: /\.vue$/,
          loader: 'vue-loader'
        },
        {
          test: /\.tsx?$/,
          loader: 'ts-loader',
          exclude: /node_modules/,
          options: {
            appendTsSuffixTo: [
              /\.vue$/
            ]
          }
        },
      ],
    },
    resolve: {
      extensions: ['.tsx', '.ts', '.js', '.vue'],
      alias: {
        vue$: 'vue/dist/vue.runtime.esm.js',
      }
    },
    plugins: [
      new VueLoaderPlugin()
    ]
  }

  if (env === 'dev') {
    cfg.mode = 'development'
    cfg.optimization = {
      usedExports: true,
    }
    return cfg
  } else {
    cfg.mode = 'production'
    return cfg
  }
}

module.exports = getConfiguration