import React from 'react';
import { connect } from 'react-redux';
import { login } from '../../actions'

function attachToButton(dispatcher) {
	return function() {
		gapi.load('auth2', function(){
			let auth2 = gapi.auth2.init({
				client_id: '900762983843-0ih1hv6b4mf4ql847ini51hhfc4svqoc.apps.googleusercontent.com',
				scope: 'email'
			});
			let btnEl = document.getElementById("googleBtn");
			if (btnEl) {
				auth2.attachClickHandler(btnEl, {},
					function(googleUser) {
						dispatcher(login(googleUser));
					}, function(error) {
						console.log(error);
					});
			}
		});
	}
}

@connect()
export default class GoogleLoginButton extends React.Component {
	static propTypes = {
		dispatch: React.PropTypes.func.isRequired
	}

	componentDidMount() {
		if (document.readyState === 'complete') {
			attachToButton(this.props.dispatch);
		} else {
			window.onload = attachToButton(this.props.dispatch);
		}
	}

	componentWillUnmount() {
		window.onload = undefined;
	}

	render() {
		return (
			<div id="googleBtn">
				<span>Login with Google</span>
			</div>
		);
	}
}
