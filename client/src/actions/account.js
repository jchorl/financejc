import { handleErrors } from './util';

export function fetchAccounts() {
    return function(dispatch) {
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        return fetch('/api/account', {
            headers,
            credentials: 'include'
        })
        .then(handleErrors)
            .then(resp => resp.json())
            .then(accounts => {
                dispatch(receiveAccountsSuccess(accounts));
                return accounts;
            })
        .catch(e => console.error(e));
    }
}

export const RECEIVE_ACCOUNTS_SUCCESS = 'RECEIVE_ACCOUNTS_SUCCESS';
function receiveAccountsSuccess(accounts) {
    return {
        type: RECEIVE_ACCOUNTS_SUCCESS,
        accounts
    }
}

export function createAccount(account) {
    return function(dispatch) {
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        headers.append('Content-Type', 'application/json');
        return fetch('/api/account', {
            headers,
            method: 'POST',
            credentials: 'include',
            body: JSON.stringify(account)
        })
        .then(handleErrors)
            .then(resp => resp.json())
            .then(account => {
                dispatch(createAccountSuccess(account));
                return account;
            })
        .catch(e => console.error(e));
    }
}

export const CREATE_ACCOUNT_SUCCESS = 'CREATE_ACCOUNT_SUCCESS';
function createAccountSuccess(account) {
    return {
        type: CREATE_ACCOUNT_SUCCESS,
        account
    }
}
