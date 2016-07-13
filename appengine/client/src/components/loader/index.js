import React from 'react';
import { render } from 'react-dom';

export default class Loader extends React.Component {
	static propTypes = {
		loading: React.PropTypes.bool.isRequired,
		children: React.PropTypes.any
	}

	render () {
		return (
			<div>
				{ this.props.loading
					? (<span>Loading</span>)
					: (
						<div>
							{ this.props.children }
						</div>
					)}
			</div>
		)
	}
}
