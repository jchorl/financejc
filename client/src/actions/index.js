export const RECEIVE_ACCOUNTS = 'RECEIVE_ACCOUNTS';
const receiveAccounts = (accounts) => {
    return {
        type: RECEIVE_ACCOUNTS,
        accounts
    };
};

export const SELECT_ACCOUNT = 'SELECT_ACCOUNT';
export const selectAccount = (id) => {
    return {
        type: SELECT_ACCOUNT,
        id
    };
};

export const ADD_ACCOUNT = 'ADD_ACCOUNT';
const addAccount = (account) => {
    return {
        type: ADD_ACCOUNT,
        account
    };
};

export const FETCHING_USER = 'FETCHING_USER';
const fetchingUser = () => {
    return {
        type: FETCHING_USER
    };
};

export const RECEIVE_USER = 'RECEIVE_USER';
const receiveUser = (user) => {
    return {
        type: RECEIVE_USER,
        user
    };
};

export const LOGOUT = 'LOGOUT';
const logoutComplete = () => {
    return {
        type: LOGOUT
    };
};

export const RECEIVE_TRANSACTIONS = 'RECEIVE_TRANSACTIONS';
const receiveTransactions = (transactions, accountId, link) => {
    return {
        type: RECEIVE_TRANSACTIONS,
        accountId,
        transactions,
        link
    };
};

export const RECEIVE_RECURRING_TRANSACTIONS = 'RECEIVE_RECURRING_TRANSACTIONS';
const receiveRecurringTransactions = (recurringTransactions, accountId) => {
    return {
        type: RECEIVE_RECURRING_TRANSACTIONS,
        accountId,
        recurringTransactions
    };
};

export const DELETE_RECURRING_TRANSACTION = 'DELETE_RECURRING_TRANSACTION';
const recurringTransactionDeleted = (id, accountId) => {
    return {
        type: DELETE_RECURRING_TRANSACTION,
        accountId,
        id
    };
};

export const RECEIVE_TEMPLATES = 'RECEIVE_TEMPLATES';
const receiveTemplates = (templates, accountId) => {
    return {
        type: RECEIVE_TEMPLATES,
        accountId,
        templates
    };
};

export const DELETE_TEMPLATE = 'DELETE_TEMPLATE';
const templateDeleted = (id, accountId) => {
    return {
        type: DELETE_TEMPLATE,
        accountId,
        id
    };
};

export const DELETE_TRANSACTION = 'DELETE_TRANSACTION';
const transactionDeleted = (id, accountId) => {
    return {
        type: DELETE_TRANSACTION,
        accountId,
        id
    };
};

export const PUT_TRANSACTION = 'PUT_TRANSACTION';
const updateTransactions = (transaction) => {
    return {
        type: PUT_TRANSACTION,
        transaction
    };
};

export const PUT_RECURRING_TRANSACTION = 'PUT_RECURRING_TRANSACTION';
const updateRecurringTransactions = (recurringTransaction) => {
    return {
        type: PUT_RECURRING_TRANSACTION,
        recurringTransaction
    };
};

export const PUT_TEMPLATE = 'PUT_TEMPLATE';
const updateTemplates = (template) => {
    return {
        type: PUT_TEMPLATE,
        template
    };
};

export const UPDATE_ACCOUNT_VALUE = 'UPDATE_ACCOUNT_VALUE';
const updateAccountValue = (accountId, delta) => {
    return {
        type: UPDATE_ACCOUNT_VALUE,
        accountId,
        delta
    };
};

export const RECEIVE_CURRENCIES = 'RECEIVE_CURRENCIES';
const receiveCurrencies = (currencies) => {
    return {
        type: RECEIVE_CURRENCIES,
        currencies
    };
};

export function fetchCurrencies() {
    return function(dispatch) {
        return fetch('/api/currencies')
            .then(response => response.json())
            .then(json => {
                dispatch(receiveCurrencies(json));
            });
    };
}

export function fetchAccounts() {
    return function(dispatch) {
        return fetch('/api/account', {
            credentials: 'include'
        })
            .then(response => response.json())
            .then(json => {
                dispatch(receiveAccounts(json));
                if (json.length) {
                    dispatch(selectAccount(json[0].id));
                }
            });
    };
}

export function fetchTransactions(accountId, next) {
    let nextStr = '';
    if (next) {
        nextStr = next;
    }

    return function(dispatch) {
        return fetch(`/api/account/${accountId}/transactions${nextStr}`, {
            credentials: 'include'
        })
            .then(response => Promise.all([response.json(), response.headers.get('Link')]))
            .then(parsed => {
                dispatch(receiveTransactions(parsed[0], accountId, parsed[1]));
            });
    };
}

export function fetchRecurringTransactions(accountId) {
    return function(dispatch) {
        return fetch(`/api/account/${accountId}/recurringTransactions`, {
            credentials: 'include'
        })
            .then(response => response.json())
            .then(parsed => {
                dispatch(receiveRecurringTransactions(parsed, accountId));
            });
    };
}

export function fetchTemplates(accountId) {
    return function(dispatch) {
        return fetch(`/api/account/${accountId}/templates`, {
            credentials: 'include'
        })
            .then(response => response.json())
            .then(parsed => {
                dispatch(receiveTemplates(parsed, accountId));
            });
    };
}

export function importData(files) {
    return function(dispatch) {
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        let data = new FormData();
        data.append('file', files[0]);
        return fetch('/api/import', {
            method: 'POST',
            credentials: 'include',
            headers: headers,
            body: data
        }).then((resp) => {
            if (resp.ok) {
                dispatch(fetchAccounts());
            }
        });
    };
}

export function fetchUser(callback) {
    return function(dispatch) {
        dispatch(fetchingUser());

        return fetch('/api/user', {
            credentials: 'include'
        })
            .then(response => {
                if (response.ok) {
                    response.json().then(json => {
                        dispatch(receiveUser(json));
                        callback && callback();
                    });
                } else {
                    dispatch(receiveUser(false));
                    callback && callback();
                }
            })
            .catch(() => {
                dispatch(receiveUser(false));
                callback && callback();
            });
    };
}

export function login(googleUser, callback) {
    return function(dispatch) {
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        headers.append('Content-Type', 'application/json');
        fetch('/api/auth', {
            method: 'POST',
            body: JSON.stringify({
                token: googleUser.getAuthResponse().id_token
            }),
            credentials: 'include',
            headers: headers
        })
            .then(response => response.json())
            .then(json => {
                dispatch(receiveUser(json));
                callback && callback();
            })
            .catch(() => dispatch(receiveUser(false)));
    };
}

export function logout(callback) {
    return function(dispatch) {
        fetch('/api/auth/logout', {
            method: 'POST',
            credentials: 'include'
        })
            .then(() => {
                dispatch(logoutComplete());
                callback && callback();
            })
            .catch(() => dispatch(logoutComplete()));
    };
}

export function newAccount(account) {
    let headers = new Headers();
    headers.append('Accept', 'application/json');
    headers.append('Content-Type', 'application/json');

    return function(dispatch) {
        return fetch('/api/account', {
            method: 'POST',
            body: JSON.stringify(account),
            credentials: 'include',
            headers: headers
        })
            .then(response => response.json())
            .then(json => {
                dispatch(addAccount(json));
                dispatch(selectAccount(json.id));
            });
    };
}

// amountDifference is the amount to change the account value by
export function putTransaction(transaction, amountDifference) {
    let headers = new Headers();
    headers.append('Accept', 'application/json');
    headers.append('Content-Type', 'application/json');

    // if editing a transaction
    if (transaction.id) {
        return function(dispatch) {
            return fetch('/api/transaction', {
                method: 'PUT',
                body: JSON.stringify(transaction),
                credentials: 'include',
                headers: headers
            })
                .then(response => response.json())
                .then(json => dispatch(updateTransactions(json)))
                .then(resp => dispatch(updateAccountValue(resp.transaction.accountId, amountDifference)));
        };
    }

    // new transaction
    return function(dispatch) {
        return fetch(`/api/account/${transaction.accountId}/transactions`, {
            method: 'POST',
            body: JSON.stringify(transaction),
            credentials: 'include',
            headers: headers
        })
            .then(response => response.json())
            .then(json => dispatch(updateTransactions(json)))
            .then(resp => dispatch(updateAccountValue(resp.transaction.accountId, amountDifference)));
    };
}

export function deleteTransaction(id, accountId, amount) {
    let headers = new Headers();
    headers.append('Accept', 'application/json');
    headers.append('Content-Type', 'application/json');

    return function(dispatch) {
        return fetch(`/api/transaction/${id}`, {
            method: 'DELETE',
            credentials: 'include',
            headers: headers
        }).then((resp) => {
            if (resp.ok) {
                dispatch(transactionDeleted(id, accountId));
                dispatch(updateAccountValue(accountId, -amount));
            }
        });
    };
}

export function putRecurringTransaction(recurringTransaction) {
    let headers = new Headers();
    headers.append('Accept', 'application/json');
    headers.append('Content-Type', 'application/json');

    // if editing a transaction
    if (recurringTransaction.id) {
        return function(dispatch) {
            return fetch('/api/recurringTransaction', {
                method: 'PUT',
                body: JSON.stringify(recurringTransaction),
                credentials: 'include',
                headers: headers
            })
                .then(response => response.json())
                .then(json => dispatch(updateRecurringTransactions(json)));
        };
    }

    return function(dispatch) {
        return fetch(`/api/account/${recurringTransaction.transaction.accountId}/recurringTransactions`, {
            method: 'POST',
            body: JSON.stringify(recurringTransaction),
            credentials: 'include',
            headers: headers
        })
            .then(response => response.json())
            .then(json => dispatch(updateRecurringTransactions(json)));
    };
}

export function putTemplate(template) {
    let headers = new Headers();
    headers.append('Accept', 'application/json');
    headers.append('Content-Type', 'application/json');

    // if editing a transaction
    if (template.id) {
        return function(dispatch) {
            return fetch('/api/template', {
                method: 'PUT',
                body: JSON.stringify(template),
                credentials: 'include',
                headers: headers
            })
                .then(response => response.json())
                .then(json => dispatch(updateTemplates(json)));
        };
    }

    return function(dispatch) {
        return fetch(`/api/account/${template.accountId}/templates`, {
            method: 'POST',
            body: JSON.stringify(template),
            credentials: 'include',
            headers: headers
        })
            .then(response => response.json())
            .then(json => dispatch(updateTemplates(json)));
    };
}

export function deleteTemplate(id, accountId) {
    let headers = new Headers();
    headers.append('Accept', 'application/json');
    headers.append('Content-Type', 'application/json');

    return function(dispatch) {
        return fetch(`/api/template/${id}`, {
            method: 'DELETE',
            credentials: 'include',
            headers: headers
        }).then((resp) => {
            if (resp.ok) {
                dispatch(templateDeleted(id, accountId));
            }
        });
    };
}

export function deleteRecurringTransaction(id, accountId) {
    let headers = new Headers();
    headers.append('Accept', 'application/json');
    headers.append('Content-Type', 'application/json');

    return function(dispatch) {
        return fetch(`/api/recurringTransaction/${id}`, {
            method: 'DELETE',
            credentials: 'include',
            headers: headers
        }).then((resp) => {
            if (resp.ok) {
                dispatch(recurringTransactionDeleted(id, accountId));
            }
        });
    };
}
