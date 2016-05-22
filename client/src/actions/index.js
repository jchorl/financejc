export const REQUEST_ACCOUNTS = 'REQUEST_ACCOUNTS';
export const requestAccounts = () => {
	return {
		type: REQUEST_ACCOUNTS
	}
}

export const RECEIVE_ACCOUNTS = 'RECEIVE_ACCOUNTS';
export const receiveAccounts = (json) => {
	return {
		type: RECEIVE_ACCOUNTS,
		accounts: json
	}
}

export function fetchAccounts() {
	return function(dispatch) {
		dispatch(requestAccounts());

		return fetch(`/account`, {
			credentials: 'include'
		})
		.then(response => response.json())
		.then(json => dispatch(receivePosts(json)))
	}
}
