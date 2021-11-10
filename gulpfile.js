const gulp = require('gulp')
const del = require('del')
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

gulp.task('critical', gulp.series('critical:index'))
gulp.task('optimize', gulp.series('critical'))
