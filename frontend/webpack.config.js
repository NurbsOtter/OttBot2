const path = require('path');

module.exports ={
	entry:'./app.js',
	output:{
		path:path.resolve('build'),
		filename:'app.bundle.js'
	},
	module:{
		loaders:[
			{test:/\.js$/, loader:'babel-loader',exclude: /node_modules/},
			{test:/\.jsx$/, loader:'babel-loader',exclude: /node_modules/}
		]
	}
}