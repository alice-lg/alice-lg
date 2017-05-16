'use strict';

/**
 * Task: bundle
 *
 * Bundle dependencies into a single file.
 *
 * See: config.bundle
 */

var gulp     = require('gulp');
var rename   = require('gulp-rename');
var concat   = require('gulp-concat');
var cssmin   = require('gulp-cssmin');
var uglify   = require('gulp-uglify');


gulp.task('bundle', function(){
  // Get bundles (js and css) from config
  var bundle = global.config.bundle;

  // Process scripts
  if ( bundle['js'] ) {
    for(var name in bundle.js) {
      var files = bundle.js[name];
      gulp.src(files)
        .pipe(concat(name + '.js'))
        .pipe(gulp.dest('build/js'))
        .pipe(uglify())
        .pipe(rename({suffix: '.min'}))
        .pipe(gulp.dest('build/js'));
    }
  }

  // Process css
  if ( bundle['css'] ) {
    for(var name in bundle.css) {
      var files = bundle.css[name];
      gulp.src(files)
        .pipe(concat(name + '.css'))
        .pipe(gulp.dest('build/css'))
        .pipe(cssmin())
        .pipe(rename({suffix: '.min'}))
        .pipe(gulp.dest('build/css'));
    }
  }

});


