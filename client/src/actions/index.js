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

export const LOGOUT = 'LOGOUT';
export const logoutComplete = () => {
  return {
    type: LOGOUT
  };
}

export const RECEIVE_TRANSACTIONS = 'RECEIVE_TRANSACTION';
export const receiveTransactions = (transactions, accountId, link) => {
  return {
    type: RECEIVE_TRANSACTIONS,
    accountId,
    transactions,
    link
  };
}

export const PUT_TRANSACTION = 'PUT_TRANSACTION';
export const updateTransactions = (transaction) => {
  return {
    type: PUT_TRANSACTION,
    transaction
  };
}

export const UPDATE_ACCOUNT_VALUE = 'UPDATE_ACCOUNT_VALUE';
export const updateAccountValue = (accountId, delta) => {
  return {
    type: UPDATE_ACCOUNT_VALUE,
    accountId,
    delta
  };
}

export const RECEIVE_CURRENCIES = 'RECEIVE_CURRENCIES';
export const receiveCurrencies = (currencies) => {
  return {
    type: RECEIVE_CURRENCIES,
    currencies
  };
}

export function fetchCurrencies() {
  return function(dispatch) {
    return fetch(`/api/currencies`)
      .then(response => response.json())
      .then(json => {
        dispatch(receiveCurrencies(json));
      });
  }
}

export function fetchAccounts() {
  return function(dispatch) {
    dispatch(requestAccounts());

    return fetch(`/api/account`, {
      credentials: 'include'
    })
      .then(response => response.json())
      .then(json => {
        dispatch(receiveAccounts(json));
      });
  }
}

export function fetchTransactions(accountId) {
  return function(dispatch) {
    return fetch(`/api/account/${accountId}/transactions`, {
      credentials: 'include'
    })
      .then(response => Promise.all([response.json(), response.headers.get('Link')]))
      .then(parsed => {
        dispatch(receiveTransactions(parsed[0], accountId, parsed[1]));
      });
  }
}

export function importData() {
  return function(dispatch) {
    let headers = new Headers();
    headers.append("Accept", "application/json");
    headers.append("Content-Type", "application/json");
    return fetch(`/api/import`, {
      method: 'POST',
      credentials: 'include',
      headers: headers
    }).then((resp) => {
      if (resp.ok) {
        dispatch(fetchAccounts());
      }
    });
  }
}

export function fetchAuth() {
  return function(dispatch) {
    dispatch(checkAuth());

    return fetch(`/api/auth`, {
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
    fetch('/api/auth', {
      method: 'POST',
      body: JSON.stringify({
        token: googleUser.getAuthResponse().id_token
      }),
      credentials: 'include',
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

export function logout() {
  return function(dispatch) {
    fetch('/api/auth/logout', {
      method: 'POST',
      credentials: 'include'
    })
      .then((resp) => {
        dispatch(logoutComplete());
      })
      .catch(() => dispatch(logoutComplete()));
  }
}

// amountDifference is the amount to change the account value by
export function putTransaction(transaction, amountDifference) {
  let headers = new Headers();
  headers.append("Accept", "application/json");
  headers.append("Content-Type", "application/json");

  // if editing a transaction
  if (transaction.id) {
    return function(dispatch) {
      return fetch(`/api/transaction`, {
        method: 'PUT',
        body: JSON.stringify(transaction),
        credentials: 'include',
        headers: headers
      })
        .then(response => response.json())
        .then(json => dispatch(updateTransactions(json)))
        .then(resp => dispatch(updateAccountValue(resp.transaction.account, amountDifference)));
    }
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
      .then(resp => dispatch(updateAccountValue(resp.transaction.account, amountDifference)));
  }
}
