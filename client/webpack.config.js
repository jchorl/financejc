var webpack = require('webpack');
var path = require('path');

var BUILD_DIR = path.resolve(__dirname, 'dest');
var APP_DIR = path.resolve(__dirname, 'src');

var config = {
	entry: APP_DIR,
	output: {
		path: BUILD_DIR,
		filename: 'bundle.js'
	},
	module: {
		loaders : [
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
			}
		]
	}
};

module.exports = config;