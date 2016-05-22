import React from 'react';
import { render } from 'react-dom';
import { connect } from 'react-redux';
import { fetchAuth } from '../../actions'
import GoogleLoginButton from '../googleLoginButton'

function mapStateToProps(state) {
	  return { auth: state.auth }
}

class App extends React.Component {
	constructor(props) {
		super(props);
		props.dispatch(fetchAuth());
	}

	render () {
		if (this.props.auth.authd) {
			return <div>Authd</div>;
		}
		return <GoogleLoginButton />;
	}
}

export default connect(mapStateToProps)(App)