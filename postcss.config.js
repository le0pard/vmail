const browserlist = require('./browserslist.config')

module.exports = {
  plugins: [
    require('postcss-import')({
      path: ['app/javascript/src']
    }),
    require('postcss-preset-env')({
      stage: 2,
      browsers: browserlist,
      features: {
        'custom-properties': {
          strict: false,
          warnings: false,
          preserve: true
        },
        'custom-media-queries': true
      }
    }),
    require('postcss-reporter')
  ]
}
