export const REQUEST_LOGIN = 'REQUEST_LOGIN';
export const requestLogin = () => {
	return {
		type: REQUEST_LOGIN
	};
}

export const REQUEST_ACCOUNTS = 'REQUEST_ACCOUNTS';
export const requestAccounts = () => {
	return {
		type: REQUEST_ACCOUNTS
	};
}

export const RECEIVE_ACCOUNTS = 'RECEIVE_ACCOUNTS';
export const receiveAccounts = (json) => {
	return {
		type: RECEIVE_ACCOUNTS,
		accounts: json
	};
}

export const CHECK_AUTH = 'CHECK_AUTH';
export const checkAuth = () => {
	return {
		type: CHECK_AUTH
	};
}

export const RECEIVE_AUTH = 'RECEIVE_AUTH';
export const receiveAuth = (authd) => {
	return {
		type: RECEIVE_AUTH,
		authd
	};
}

export const UPDATE_TRANSACTION = 'UPDATE_TRANSACTION';
export const updateTransaction = (transaction) => {
	return {
		type: updateTransaction,
		transaction
	};
}

export function fetchAccounts() {
	return function(dispatch) {
		dispatch(requestAccounts());

		return fetch(`/account`, {
			credentials: 'include'
		})
		.then(response => response.json())
		.then(json => dispatch(receiveAccounts(json)));
	}
}

export function fetchAuth() {
	return function(dispatch) {
		dispatch(checkAuth());

		return fetch(`/auth`, {
			credentials: 'include'
		})
		.then((resp) => {
			if (resp.ok) {
				dispatch(receiveAuth(true));
			} else {
				dispatch(receiveAuth(false));
			}
		})
		.catch(() => dispatch(receiveAuth(false)));
	}
}

export function login(googleUser) {
	return function(dispatch) {
		let headers = new Headers();
		headers.append("Accept", "application/json");
		headers.append("Content-Type", "application/json");
		fetch('/auth', {
			method: 'POST',
			body: JSON.stringify({
				token: googleUser.getAuthResponse().id_token
			}),
			credentials: 'same-origin',
			headers: headers
		})
		.then((resp) => {
			if (resp.ok) {
				dispatch(receiveAuth(true));
			} else {
				dispatch(receiveAuth(false));
			}
		})
		.catch(() => dispatch(receiveAuth(false)));
	}
}

export function putTransaction(transaction) {
	return function(dispatch) {
		let headers = new Headers();
		headers.append("Accept", "application/json");
		headers.append("Content-Type", "application/json");
		return fetch(`/transaction/${transaction.id}`, {
			method: 'PUT',
			body: JSON.stringify(transaction),
			credentials: 'include',
			headers: headers
		})
		.then(response => response.json())
		.then(json => dispatch(updateTransaction(transaction)));
	}
}
