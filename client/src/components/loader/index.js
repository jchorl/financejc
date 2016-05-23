import React from 'react';
import { render } from 'react-dom';

export default class Loader extends React.Component {
	static propTypes = {
		loading: React.PropTypes.bool.isRequired
	}

	render () {
		if (this.props.loading) {
			return <span>Loading</span>
		}
		return <div>{this.props.children}</div>;
	}
}
