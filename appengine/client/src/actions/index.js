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

export const PUT_TRANSACTION = 'PUT_TRANSACTION';
export const updateTransactions = (transaction) => {
	return {
		type: PUT_TRANSACTION,
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
			let accounts = normalized.entities.accounts || {};
			let transactions = normalized.entities.transactions || {};
			dispatch(receiveAccounts(accounts));
			dispatch(receiveTransactions(transactions));
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

export function putTransaction(transaction) {
	let headers = new Headers();
	headers.append("Accept", "application/json");
	headers.append("Content-Type", "application/json");

	// if editing a transaction
	if (transaction.id) {
		return function(dispatch) {
			return fetch(`/transaction/${transaction.id}`, {
				method: 'PUT',
				body: JSON.stringify(transaction),
				credentials: 'include',
				headers: headers
			})
			.then(response => response.json())
			.then(json => dispatch(updateTransactions(json)));
		}
	}

	// new transaction
	return function(dispatch) {
		return fetch(`/account/${transaction.accountId}/transactions`, {
			method: 'POST',
			body: JSON.stringify(transaction),
			credentials: 'include',
			headers: headers
		})
		.then(response => response.json())
		.then(json => dispatch(updateTransactions(json)));
	}
}