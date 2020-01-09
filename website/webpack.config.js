/* eslint-disable */
const path = require('path')
const VueLoaderPlugin = require('vue-loader/lib/plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const { CleanWebpackPlugin } = require('clean-webpack-plugin');

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
        {
          test: /\.scss$/,
          use: [
            'vue-style-loader',
            'css-loader',
            {
              loader: 'sass-loader',
              options: {
                prependData: '@import "src/assets/scss/all";',
              },
            },
          ]
        },
      ],
    },
    resolve: {
      extensions: ['.tsx', '.ts', '.js', '.vue'],
      alias: {
        vue$: 'vue/dist/vue.runtime.esm.js',
        '@': path.resolve("./src"),
        '@scss': path.resolve("./src/assets/scss")
      }
    },
    plugins: [
      new CleanWebpackPlugin(),
      new VueLoaderPlugin(),
      new HtmlWebpackPlugin({
        title: 'Dionysos',
        template: 'src/index.html'
      }),
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