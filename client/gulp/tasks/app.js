'use strict';

/**
 * Task: app
 *
 * Compile the main app
 */

var gulp       = require('gulp');
var uglify     = require('gulp-uglify');
var rename     = require('gulp-rename');
var size       = require('gulp-size');

var browserify = require('browserify');
var babelify   = require('babelify');

var source     = require('vinyl-source-stream');

// == Register task: app 
gulp.task('app', function(){
  var entries = ['./app.jsx'];

  if (process.env.DISABLE_LOGGING) {
    entries.unshift('./no_log.jsx');
  }

  var bundler = browserify({
    entries: entries,
    extensions: ['.jsx'],
    paths: ['./node_modules', './']
  });


  return bundler.transform(
      babelify.configure({
        presets: ["es2015", "react"]
      })
    )
    .bundle()
    .pipe(source('app.js'))
    .pipe(gulp.dest('build/js'))
    .pipe(size());

});

gulp.task('appmin', ['app'], function() {
  return gulp.src('build/js/app.js')
    .pipe(uglify())
    .pipe(size())
    .pipe(rename({suffix: '.min'}))
    .pipe(gulp.dest('build/js'));
});


