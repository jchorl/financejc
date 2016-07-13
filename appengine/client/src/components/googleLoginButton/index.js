import React from 'react';
import { connect } from 'react-redux';
import { login } from '../../actions'

@connect()
export default class GoogleLoginButton extends React.Component {
	constructor(props) {
		super(props);
		window.onload = function() {
			gapi.load('auth2', function(){
				let auth2 = gapi.auth2.init({
					client_id: '900762983843-0ih1hv6b4mf4ql847ini51hhfc4svqoc.apps.googleusercontent.com',
					scope: 'email'
				});
				let btnEl = document.getElementById('googleBtn');
				if (btnEl) {
					auth2.attachClickHandler(btnEl, {},
											 function(googleUser) {
												 props.dispatch(login(googleUser));
											 }, function(error) {
												 console.log(error);
											 });
				}
			});
		};
	}

	render () {
		return (
			<div id="googleBtn">
				<span>Login with Google</span>
			</div>
		);
	}
}
