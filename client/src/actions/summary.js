import { handleErrors } from './util';

export function fetchSummary(since) {
    return function(dispatch) {
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        return fetch(`/api/summary?since=${since.toISOString()}`, {
            headers,
            credentials: 'include'
        })
        .then(handleErrors)
            .then(resp => resp.json())
            .then(transactions => {
                dispatch(receiveSummarySuccess(transactions));
                return transactions;
            })
        .catch(e => console.error(e));
    }
}

export const RECEIVE_SUMMARY_SUCCESS = 'RECEIVE_SUMMARY_SUCCESS';
function receiveSummarySuccess(transactions) {
    return {
        type: RECEIVE_SUMMARY_SUCCESS,
        transactions
    }
}
