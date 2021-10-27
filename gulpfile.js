const gulp = require('gulp')
const del = require('del')
const purgecss = require('gulp-purgecss')
const gzip = require('gulp-gzip')
const critical = require('critical').stream

const criticalOptions = {
  base: 'build/',
  inline: true,
  width: 1440,
  height: 1024
}

// Assetts cleanup
gulp.task('cleanup:assets', () => {
  return del([
    '.tmp/dist/**/*'
  ])
})

// Purgecss for app.css
gulp.task('purgecss:app_css', () => {
  return gulp.src('build/app-*.css')
    .pipe(purgecss({
      content: ['build/**/*.html', 'webpack/controllers/**/*.js'],
      safelist: {
        greedy: [/plyr/]
      }
    }))
    .pipe(gulp.dest('build'))
})

// Gzip app.css after purgecss
gulp.task('purgecss:recompress_app_css', () => {
  return gulp.src('build/app-*.css')
    .pipe(gzip({
      append: true,
      threshold: '10kb',
      gzipOptions: {level: 9, memLevel: 8},
      skipGrowingFiles: true
    }))
    .pipe(gulp.dest('build'))
})

// Generate & Inline Critical-path CSS
gulp.task('critical:index', () => {
  return gulp
    .src(['build/*.html', '!build/404.html'])
    .pipe(critical(criticalOptions))
    .on('error', (err) => {
      // eslint-disable-next-line no-console
      console.error(JSON.stringify(err))
    })
    .pipe(gulp.dest('build'))
})

gulp.task('purgecss', gulp.series('purgecss:app_css', 'purgecss:recompress_app_css'))
gulp.task('critical', gulp.series('critical:index'))
gulp.task('optimize', gulp.series('purgecss', 'critical'))
