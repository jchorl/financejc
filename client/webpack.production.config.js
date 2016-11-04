var webpack = require('webpack');
var path = require('path');
var HtmlWebpackPlugin = require('html-webpack-plugin')

var BUILD_DIR = path.resolve(__dirname, 'dest');
var APP_DIR = path.resolve(__dirname, 'src');

var constants = require('./src/components/index/constants');

var config = {
  devtool: 'source-map',
  entry: APP_DIR,
  output: {
    path: BUILD_DIR,
    filename: 'bundle.js'
  },
  module: {
    loaders: [
      {
        test : /\.jsx?/,
        include : APP_DIR,
        loaders : ['babel?cacheDirectory=true,presets[]=es2015,presets[]=react,presets[]=stage-1,plugins[]=transform-decorators-legacy', 'eslint-loader'],
      },
      {
        test : /\.css$/,
        loader: 'style-loader!css-loader?modules&importLoaders=1&localIdentName=[name]_[local]_[hash:base64:5]!postcss-loader'
      }
    ]
  },
  eslint: {
    configFile: './.eslintrc'
  },
  postcss: [
    require('postcss-constants')({
      defaults: constants
    }),
    require('autoprefixer'),
    require('precss'),
    require('postcss-nested'),
    require('lost')
  ],
  plugins: [
    new webpack.DefinePlugin({
      'process.env': {
        'NODE_ENV': JSON.stringify(process.env.NODE_ENV)
      }
    }),
    new webpack.optimize.DedupePlugin(),
    new webpack.optimize.UglifyJsPlugin({
      compress: {
        screw_ie8: true
      }
    }),
    new webpack.optimize.AggressiveMergingPlugin(),
    new HtmlWebpackPlugin({
      template: APP_DIR + '/index.html',
      filename: 'index.html',
      minify: {},
      inject: 'body'
    })
  ]
};

module.exports = config;
