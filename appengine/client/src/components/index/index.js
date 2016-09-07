import React from 'react';
import { render } from 'react-dom';
import { connect } from 'react-redux';
import { fetchAuth } from '../../actions';
import AccountsPage from '../accountsPage';
import GoogleLoginButton from '../googleLoginButton';
import Loader from '../loader';
import styles from './styles.css';

@connect((state) => {
	return { auth: state.auth }
})
export default class App extends React.Component {
	static propTypes = {
		auth: React.PropTypes.object.isRequired
	}

	constructor(props) {
		super(props);
		props.dispatch(fetchAuth());
	}

	render () {
		return (
			<Loader loading={ !this.props.auth.fetched }>
				{ this.props.auth.authd ? <AccountsPage /> : <GoogleLoginButton /> }
			</Loader>
		)
	}
}