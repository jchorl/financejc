window.onload = function() {
	gapi.load('auth2', function(){
		// Retrieve the singleton for the GoogleAuth library and set up the client.
		auth2 = gapi.auth2.init({
			client_id: '900762983843-0ih1hv6b4mf4ql847ini51hhfc4svqoc.apps.googleusercontent.com',
			scope: 'email'
		});
		auth2.attachClickHandler(document.getElementById('googleBtn'), {},
								 function(googleUser) {
									 console.log('fetching');
									 fetch('/auth', {
										 method: 'POST',
										 body: JSON.stringify({
											 token: googleUser.getAuthResponse().id_token
										 }),
										 credentials: 'same-origin'
									 });
								 }, function(error) {
									 alert(JSON.stringify(error, undefined, 2));
								 });
	});
};
