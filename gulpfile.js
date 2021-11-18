const gulp = require('gulp')
const del = require('del')

// Assetts cleanup
gulp.task('cleanup:assets', () => {
  return del([
    '.tmp/dist/**/*'
  ])
})

