import React from 'react';

export default class GoogleLoginButton extends React.Component {
	constructor(props) {
		super(props);
		window.onload = function() {
			gapi.load('auth2', function(){
				let auth2 = gapi.auth2.init({
					client_id: '900762983843-0ih1hv6b4mf4ql847ini51hhfc4svqoc.apps.googleusercontent.com',
					scope: 'email'
				});
				auth2.attachClickHandler(document.getElementById('googleBtn'), {},
										 function(googleUser) {
											 var headers = new Headers();
											 headers.append("Accept", "application/json");
											 headers.append("Content-Type", "application/json");
											 fetch('/auth', {
												 method: 'POST',
												 body: JSON.stringify({
													 token: googleUser.getAuthResponse().id_token
												 }),
												 credentials: 'same-origin',
												 headers: headers
											 });
										 }, function(error) {
											 alert(JSON.stringify(error, undefined, 2));
										 });
			});
		};
	}
	render () {
		return (
			<div id="googleBtn">
				<span>Google</span>
			</div>
		)
	}
}
