'use strict'

var path = require('path');

module.exports = {
  bail: true,
  cache: true,
  context: path.join(__dirname, 'src'),
  debug: true,
  devtool: '#inline-source-map',
  entry: {
    javascript: './js/main.js',
    html: './index.html'
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
        loader: "url-loader?mimetype=image/png",
        test: /\.png$/,
      },
      {
        loader: "file-loader",
        test: /\.(ttf|eot|svg)(\?v=[0-9]\.[0-9]\.[0-9])?$/
      },
      {
        loader: "url-loader",
        test: /\.woff(2)?(\?v=\d\.\d\.\d)?$/,
        query: {
          limit: 1000,
          mimetype: 'application/font-woff'
        }
      }
    ]
  },
  output: {
    filename: 'bundle.js',
    path: path.join(__dirname, 'build'),
    publicPath: './',
    sourcePrefix: ' '
  },
  resolve: {
    extensions: [ '', '.js' ]
  },
  stats: {
    colors: true,
    reasons: true
  }
}
