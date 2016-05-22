import React from 'react';
import { render } from 'react-dom';
import { connect } from 'react-redux';
import { fetchAuth } from '../../actions';
import AccountsPage from '../accountsPage';
import GoogleLoginButton from '../googleLoginButton';
import Loader from '../loader';

@connect((state) => {
	return { auth: state.auth }
})
export default class App extends React.Component {
	constructor(props) {
		super(props);
		props.dispatch(fetchAuth());
	}

	render () {
		return (
			<Loader loading={ this.props.auth.isFetching }>
				{ this.props.auth.authd ? <AccountsPage /> : <GoogleLoginButton /> }
			</Loader>
		)
	}
}