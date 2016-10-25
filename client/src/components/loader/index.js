import React from 'react';
import { render } from 'react-dom';
import styles from './loader.css';

export default class Loader extends React.Component {
	static propTypes = {
		loading: React.PropTypes.bool.isRequired,
		children: React.PropTypes.any
	}

	render () {
		return (
			<div>
				{ this.props.loading
					? (<div className={ styles.loader }>Loading...</div>)
					: (
						<div>
							{ this.props.children }
						</div>
					)}
			</div>
		)
	}
}
