'use strict'

var path = require('path')

module.exports = {
  bail: false,
  cache: true,
  context: path.join(__dirname, 'src'),
  debug: true,
  devtool: '#inline-source-map',
  entry: {
    javascript: './scripts/Main.jsx',
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
        loader: "url-loader",
        test: /\.woff(2)?(\?v=\d\.\d\.\d)?$/,
        query: {
          limit: 1000,
          mimetype: 'application/font-woff'
        }
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
  resolve: {
    extensions: [ '', '.js', '.jsx' ]
  },
  stats: {
    colors: true,
    reasons: true
  }
}
