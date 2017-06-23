import { handleErrors } from './util';

export function fetchCurrencies() {
    return function(dispatch) {
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        return fetch('/api/currencies', {
            headers
        })
        .then(handleErrors)
            .then(resp => resp.json())
            .then(currencies => {
                dispatch(receiveCurrenciesSuccess(currencies));
                return currencies;
            })
        .catch(e => console.error(e));
    }
}

export const RECEIVE_CURRENCIES_SUCCESS = 'RECEIVE_CURRENCIES_SUCCESS';
function receiveCurrenciesSuccess(currencies) {
    return {
        type: RECEIVE_CURRENCIES_SUCCESS,
        currencies
    }
}
