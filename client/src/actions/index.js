import { normalize, Schema, arrayOf } from 'normalizr';

const accountSchema = new Schema('accounts');
const transactionSchema = new Schema('transactions');

accountSchema.define({
	transactions: arrayOf(transactionSchema)
});

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
export const receiveAccounts = (accounts) => {
	return {
		type: RECEIVE_ACCOUNTS,
		accounts
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

export const RECEIVE_TRANSACTIONS = 'RECEIVE_TRANSACTION';
export const receiveTransactions = (transactions) => {
	return {
		type: RECEIVE_TRANSACTIONS,
		transactions
	};
}

export const UPDATE_TRANSACTION = 'UPDATE_TRANSACTION';
export const updateTransaction = (transaction) => {
	return {
		type: UPDATE_TRANSACTION,
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
		.then(json => {
			let normalized = normalize(json, arrayOf(accountSchema));
			dispatch(receiveAccounts(normalized.entities.accounts))
			dispatch(receiveTransactions(normalized.entities.transactions))
		});
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

export function putTransaction(index, transaction) {
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
		.then(json => dispatch(updateTransaction(index, transaction)));
	}
}
