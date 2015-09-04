'use strict'

var autoprefixer = require('autoprefixer-core');
var path = require('path');
var webpack = require('webpack');
var HtmlWebpackPlugin = require('html-webpack-plugin')

module.exports = {
  bail: false,
  cache: false,
  context: path.join(__dirname, 'src'),
  debug: false,
  entry: {
    javascript: './scripts/Main.jsx',
    html: './index.html'
  },
  eslint: {
    failOnError: true
  },
  module: {
    loaders: [
      {
        loader: 'style-loader!css-loader!postcss-loader',
        test: /\.css$/
      },
      {
        loader: 'style-loader!css-loader!postcss-loader!less-loader',
        test: /\.less$/
      },
      {
        loader: "url-loader?mimetype=image/gif",
        test: /\.gif$/,
      },
      {
        loader: 'file?name=[name].[ext]',
        test: /\.html$/
      },
      {
        exclude: /node_modules/,
        loaders: [ 'babel-loader' ],
        test: /\.jsx?$/
      },
      {
        loader: "url-loader?mimetype=image/png",
        test: /\.png$/,
      },
      {
        exclude: /node_modules/,
        loader: require.resolve('./svg.loader'),
        test: /\.svg(\?t=custom)$/
      },
      {
        loader: "file-loader",
        test: /\.(ttf|eot|svg)(\?v=[0-9]\.[0-9]\.[0-9])?$/
      },
      {
        loader: "url-loader?limit=10000&minetype=application/font-woff",
        test: /\.woff(2)?(\?v=\d\.\d\.\d)?$/
      }
    ],
    preloaders: [
      {
        exclude: /node_modules/,
        loader: 'eslint-loader',
        test: /\.jsx?$/
      }
    ]
  },
  output: {
    filename: 'bundle.js',
    path: path.join(__dirname, 'build'),
    publicPath: './',
    sourcePrefix: ' '
  },
  plugins: [
    new HtmlWebpackPlugin({
      segmentKey: 'New9ADhDv3J4ipQ3OrfsoWo9DMlVCIxO',
      template: './src/index.html'
    }),
    new webpack.optimize.UglifyJsPlugin({
      compress: {
        warnings: false
      },
      mangle: true,
      minimize: true,
      sourceMap: false
    })
  ],
  postcss: [
    autoprefixer({
      browsers: ['last 2 versions']
    })
  ],
  resolve: {
    extensions: [ '', '.js', '.jsx' ]
  },
  stats: {
    colors: true,
    reasons: true
  }
}
