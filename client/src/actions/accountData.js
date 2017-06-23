import { handleErrors } from './util';

export function fetchTransactions(accountId) {
    return function(dispatch) {
        dispatch(fetchTransactionsStart(accountId));
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        return fetch(`/api/account/${accountId}/transactions`, {
            headers,
            credentials: 'include'
        })
        .then(handleErrors)
            .then(resp => resp.json())
            .then(transactions => {
                dispatch(receiveTransactionsSuccess(accountId, transactions));
                return transactions;
            })
        .catch(e => console.error(e));
    }
}

export const FETCH_TRANSACTIONS_START = 'FETCH_TRANSACTIONS_START';
function fetchTransactionsStart(accountId) {
    return {
        type: FETCH_TRANSACTIONS_START,
        accountId
    }
}

export const RECEIVE_TRANSACTIONS_SUCCESS = 'RECEIVE_TRANSACTIONS_SUCCESS';
function receiveTransactionsSuccess(accountId, transactions) {
    return {
        type: RECEIVE_TRANSACTIONS_SUCCESS,
        accountId,
        transactions
    }
}

export function newTransaction(transaction, amountDifference) {
    return function(dispatch) {
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        headers.append('Content-Type', 'application/json');
        return fetch(`/api/account/${transaction.accountId}/transactions`, {
            headers,
            method: 'POST',
            credentials: 'include',
            body: JSON.stringify(transaction)
        })
        .then(handleErrors)
            .then(resp => resp.json())
            .then(transaction => {
                dispatch(editTransactionSuccess(transaction, amountDifference));
                return transaction;
            })
        .catch(e => console.error(e));
    }
}

export function editTransaction(transaction, amountDifference) {
    return function(dispatch) {
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        headers.append('Content-Type', 'application/json');
        return fetch(`/api/transaction`, {
            headers,
            method: 'PUT',
            credentials: 'include',
            body: JSON.stringify(transaction)
        })
        .then(handleErrors)
            .then(resp => resp.json())
            .then(transaction => {
                dispatch(editTransactionSuccess(transaction, amountDifference));
                return transaction;
            })
        .catch(e => console.error(e));
    }
}

export const EDIT_TRANSACTION_SUCCESS = 'EDIT_TRANSACTION_SUCCESS';
function editTransactionSuccess(transaction, amountDifference) {
    return {
        type: EDIT_TRANSACTION_SUCCESS,
        transaction,
        amountDifference
    }
}

export function fetchTemplates(accountId) {
    return function(dispatch) {
        dispatch(fetchTemplatesStart(accountId));
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        return fetch(`/api/account/${accountId}/templates`, {
            headers,
            credentials: 'include'
        })
        .then(handleErrors)
            .then(resp => resp.json())
            .then(templates => {
                dispatch(receiveTemplatesSuccess(accountId, templates));
                return templates;
            })
        .catch(e => console.error(e));
    }
}

export const FETCH_TEMPLATES_START = 'FETCH_TEMPLATES_START';
function fetchTemplatesStart(accountId) {
    return {
        type: FETCH_TEMPLATES_START,
        accountId
    }
}

export const RECEIVE_TEMPLATES_SUCCESS = 'RECEIVE_TEMPLATES_SUCCESS';
function receiveTemplatesSuccess(accountId, templates) {
    return {
        type: RECEIVE_TEMPLATES_SUCCESS,
        accountId,
        templates
    }
}

export function editTemplate(template) {
    return function(dispatch) {
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        headers.append('Content-Type', 'application/json');
        return fetch(`/api/template`, {
            headers,
            method: 'PUT',
            credentials: 'include',
            body: JSON.stringify(template)
        })
        .then(handleErrors)
            .then(resp => resp.json())
            .then(template => {
                dispatch(editTemplateSuccess(template));
                return template;
            })
        .catch(e => console.error(e));
    }
}

export function newTemplate(template) {
    return function(dispatch) {
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        headers.append('Content-Type', 'application/json');
        return fetch(`/api/account/${template.accountId}/templates`, {
            headers,
            method: 'POST',
            credentials: 'include',
            body: JSON.stringify(template)
        })
        .then(handleErrors)
            .then(resp => resp.json())
            .then(template => {
                dispatch(editTemplateSuccess(template));
                return template;
            })
        .catch(e => console.error(e));
    }
}

export const EDIT_TEMPLATE_SUCCESS = 'EDIT_TEMPLATE_SUCCESS';
function editTemplateSuccess(template) {
    return {
        type: EDIT_TEMPLATE_SUCCESS,
        template
    }
}

export function fetchRecurringTransactions(accountId) {
    return function(dispatch) {
        dispatch(fetchRecurringTransactionsStart(accountId));
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        return fetch(`/api/account/${accountId}/recurringTransactions`, {
            headers,
            credentials: 'include'
        })
        .then(handleErrors)
            .then(resp => resp.json())
            .then(recurringTransactions => {
                dispatch(receiveRecurringTransactionsSuccess(accountId, recurringTransactions));
                return recurringTransactions;
            })
        .catch(e => console.error(e));
    }
}

export const FETCH_RECURRING_TRANSACTIONS_START = 'FETCH_RECURRING_TRANSACTIONS_START';
function fetchRecurringTransactionsStart(accountId) {
    return {
        type: FETCH_RECURRING_TRANSACTIONS_START,
        accountId
    }
}

export const RECEIVE_RECURRING_TRANSACTIONS_SUCCESS = 'RECEIVE_RECURRING_TRANSACTIONS_SUCCESS';
function receiveRecurringTransactionsSuccess(accountId, recurringTransactions) {
    return {
        type: RECEIVE_RECURRING_TRANSACTIONS_SUCCESS,
        accountId,
        recurringTransactions
    }
}

export function editRecurringTransaction(recurringTransaction) {
    return function(dispatch) {
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        headers.append('Content-Type', 'application/json');
        return fetch(`/api/recurringTransaction`, {
            headers,
            method: 'PUT',
            credentials: 'include',
            body: JSON.stringify(recurringTransaction)
        })
        .then(handleErrors)
            .then(resp => resp.json())
            .then(recurringTransaction => {
                dispatch(editRecurringTransactionSuccess(recurringTransaction));
                return recurringTransaction;
            })
        .catch(e => console.error(e));
    }
}

export function newRecurringTransaction(recurringTransaction) {
    return function(dispatch) {
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        headers.append('Content-Type', 'application/json');
        return fetch(`/api/account/${recurringTransaction.transaction.accountId}/recurringTransactions`, {
            headers,
            method: 'POST',
            credentials: 'include',
            body: JSON.stringify(recurringTransaction)
        })
        .then(handleErrors)
            .then(resp => resp.json())
            .then(recurringTransaction => {
                dispatch(editRecurringTransactionSuccess(recurringTransaction));
                return recurringTransaction;
            })
        .catch(e => console.error(e));
    }
}

export const EDIT_RECURRING_TRANSACTION_SUCCESS = 'EDIT_RECURRING_TRANSACTION_SUCCESS';
function editRecurringTransactionSuccess(recurringTransaction) {
    return {
        type: EDIT_RECURRING_TRANSACTION_SUCCESS,
        recurringTransaction
    }
}
