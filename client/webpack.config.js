var webpack = require('webpack');
var path = require('path');

var BUILD_DIR = path.resolve(__dirname, 'dest');
var APP_DIR = path.resolve(__dirname, 'src');

var constants = require('./src/components/index/constants');

var config = {
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
				loader : 'babel',
				query: {
					cacheDirectory: true,
					plugins: [
						'transform-decorators-legacy'
					],
					presets: ['es2015', 'react', 'stage-1']
				}
			},
			{
				test : /\.css$/,
				loader: 'style-loader!css-loader?modules&importLoaders=1&localIdentName=[name]_[local]_[hash:base64:5]!postcss-loader'
			}
		]
	},
	postcss: [
		require('postcss-constants')({
			defaults: constants
		}),
		require('autoprefixer'),
		require('precss'),
		require('postcss-nested'),
		require('lost')
	]
};

module.exports = config;