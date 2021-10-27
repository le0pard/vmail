// Example webpack configuration with asset fingerprinting in production.

const path = require('path')
const webpack = require('webpack')
const TerserPlugin = require('terser-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin')
const WorkboxPlugin = require('workbox-webpack-plugin')
const WebpackAssetsManifest = require('webpack-assets-manifest')

const browserList = require('./browserslist.config')

// set NODE_ENV=production on the environment to add asset fingerprints
const currentEnv = process.env.NODE_ENV || 'development'
const isProduction = currentEnv === 'production'

const preScripts = {
  development: [],
  production: []
}

const preScriptsEnv = isProduction ?
  preScripts['production'] :
  preScripts['development']

const cssLoaders = [
  MiniCssExtractPlugin.loader,
  {
    loader: 'css-loader',
    options: {
      modules: false,
      sourceMap: true
    }
  },
  {
    loader: 'postcss-loader',
    options: {
      sourceMap: true,
      postcssOptions: () => {
        const plugins = [
          ['postcss-import'],
          ['postcss-preset-env', {
            stage: 1,
            browsers: browserList,
            features: {
              'custom-properties': {
                strict: false,
                warnings: false,
                preserve: true
              }
            }
          }],
          ['lost', {
            flexbox: 'flex'
          }],
          ['rucksack-css'],
          ['postcss-browser-reporter'],
          ['postcss-reporter']
        ]

        return {plugins}
      }
    }
  },
  {
    loader: 'sass-loader',
    options: {
      sourceMap: true,
      webpackImporter: true,
      implementation: require('sass'),
      sassOptions: {
        fiber: false,
        includePaths: [path.join(__dirname, 'webpack', 'css')]
      }
    }
  }
]

let config = {
  target: 'web',
  mode: currentEnv,
  performance: {
    hints: 'warning',
    maxEntrypointSize: 1024000,
    maxAssetSize: 1024000
  },
  entry: {
    'app': preScriptsEnv.concat(['./webpack/app.js'])
  },

  output: {
    // Build assets directly in to public/webpack/, let webpack know
    // that all webpacked assets start with webpack/

    // must match config.webpack.output_dir
    path: path.join(__dirname, '.tmp', 'dist'),
    publicPath: '/',
    filename: isProduction ? '[name]-[contenthash].js' : '[name].js',
    assetModuleFilename: 'assets/[name]-[contenthash][ext]'
  },

  resolve: {
    modules: [
      path.join(__dirname, 'webpack'),
      path.join(__dirname, 'source/images'),
      path.join(__dirname, 'node_modules')
    ],
    extensions: ['.js', '.json']
  },

  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: [
          'babel-loader'
        ]
      },
      {
        test: /\.(gif|jpg|png|woff|woff2|eot|ttf|svg|ico)$/,
        type: 'asset',
        parser: {
          dataUrlCondition: {
            maxSize: 10000
          }
        }
      },
      {
        test: /\.(scss|sass)$/,
        sideEffects: true,
        use: cssLoaders
      }
    ]
  },

  plugins: [
    new MiniCssExtractPlugin({
      filename: isProduction ? '[name]-[contenthash].css' : '[name].css'
    })
  ],

  optimization: {
    sideEffects: false,
    usedExports: true
  }
}

config.optimization = config.optimization || {}
config.optimization.runtimeChunk = false
config.optimization.splitChunks = {
  chunks: 'async'
}

if (isProduction) {
  config.plugins.push(
    new webpack.DefinePlugin({
      'process.env': {NODE_ENV: JSON.stringify('production')}
    })
  )
  config.optimization.minimizer = [
    new TerserPlugin({
      parallel: true,
      terserOptions: {
        parse: {
          // Let terser parse ecma 8 code but always output
          // ES6+ compliant code for older browsers
          ecma: 8
        },
        compress: {
          ecma: 2016,
          warnings: false,
          comparisons: false
        },
        mangle: {
          toplevel: true,
          safari10: false
        },
        output: {
          ecma: 2016,
          comments: false,
          ascii_only: true
        }
      }
    }),
    new CssMinimizerPlugin({
      minimizerOptions: {
        preset: [
          'default',
          {
            discardComments: {removeAll: true},
            mergeRules: false, // right now it break focus-visible
            svgo: false
          }
        ]
      }
    })
  ]
  // Source maps
  config.devtool = 'source-map'
} else {
  config.optimization.moduleIds = 'named'
  // Source maps
  config.devtool = 'inline-source-map'
}

config.plugins.push(
  new WebpackAssetsManifest({
    output: 'assets-manifest.json',
    publicPath: config.output.publicPath,
    writeToDisk: true,
    integrity: true,
    integrityHashes: ['sha256']
  }),
  new WorkboxPlugin.InjectManifest({
    swSrc: './webpack/sw.js',
    swDest: 'sw.js',
    compileSrc: true,
    maximumFileSizeToCacheInBytes: (isProduction ? 2097152 : 15730000)
  })
)

module.exports = config
