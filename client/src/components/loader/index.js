import PropTypes from 'prop-types';
import React from 'react';
import styles from './loader.css';

export default class Loader extends React.Component {
    static propTypes = {
        loading: PropTypes.bool.isRequired,
        children: PropTypes.any
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
        );
    }
}
