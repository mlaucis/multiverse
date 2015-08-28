'use strict'

var gulp = require('gulp')
var http = require('http')
var nodeStatic = require('node-static')
var plugins = require('gulp-load-plugins')()
var webpack = require('webpack')

gulp.task('build', [ 'bundle', 'favicon' ])

gulp.task('bundle', function(cb) {
  var config = require('./webpack.production.js')
  var bundler = webpack(config)

  bundler.run(reportBundle(cb))
})

gulp.task('bundle:watch', function(cb) {
  var config = require('./webpack.config.js')
  var bundler = webpack(config)

  bundler.watch(200, reportBundle(cb))
})

gulp.task('favicon', function () {
  gulp
    .src('./src/icons/favicon.png')
    .pipe(gulp.dest('./build'))
})

gulp.task('lint', function() {
  return gulp
    .src([ 'src/**/*.js*' ])
    .pipe(plugins.eslint())
    .pipe(plugins.eslint.format())
    .pipe(plugins.eslint.failOnError())
});

gulp.task('serve', [ 'bundle:watch' ], function() {
  var file = new nodeStatic.Server('./build', {
    cache: false
  })

  http.createServer(function (req, res) {
    file.serve(req, res).addListener('error', function (err) {
      file.serveFile('/index.html', 200, {}, req, res)
    })
  }).listen(8082)
})

gulp.task('sync', [ 'serve' ], function(cb) {
  var browserSync = require('browser-sync')

  browserSync({
    https: false,
    logPrefix: 'RSX',
    notify: false,
    online: false,
    open: false,
    port: 8080,
    proxy: 'localhost:8082',
    reloadDebounce: 2000,
    ui: {
      port: 8081
    }
  }, cb)

  process.on('exit', function() {
    browserSync.exit()
  })

  gulp.watch('build/**/*.*', browserSync.reload)
})

function reportBundle(cb) {
  var invoked = false

  return function(err, stats) {
    if (err) {
      throw new plugins.util.PluginError('webpack', err)
    }

    console.log(stats.toString({
      cached: true,
      cachedAssets: true,
      chunks: true,
      chunkModules: true,
      colors: plugins.util.colors.supportsColor,
      hash: true,
      timings: true,
      version: true
    }))

    if (!invoked) {
      invoked = true
      cb()
    }
  }
}
