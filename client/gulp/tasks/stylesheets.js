'use strict';

/**
 * Task: stylesheets 
 *
 * Compile stylesheets from less source.
 */

var gulp         = require('gulp');
var sass         = require('gulp-sass');
var cssmin       = require('gulp-cssmin');
var rename       = require('gulp-rename');
var autoprefixer = require('gulp-autoprefixer');



// == Register task: stylesheets 
gulp.task('stylesheets', function(){
  // Compile less files
  gulp.src('assets/scss/*.scss')
    .pipe(sass().on('error', sass.logError))
    .pipe(autoprefixer({
      browsers: ['last 2 versions'],
      cascade: false
     }))
    .pipe(gulp.dest('build/css/'))
    .pipe(cssmin())
    .pipe(rename({suffix: '.min'}))
    .pipe(gulp.dest('build/css/'));
});

