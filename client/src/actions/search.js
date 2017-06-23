import { handleErrors } from './util';

export function search(term) {
    return function(dispatch) {
        let timestamp = new Date();
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        return fetch(`/api/search?value=${term}`, {
            headers,
            credentials: 'include'
        })
        .then(handleErrors)
        .then(resp => resp.json())
        .then(results => {
            dispatch(receiveSearchResultsSuccess(timestamp, results));
            return results;
        })
        .catch(e => console.error(e));
    }
}

export const CLEAR_SEARCH_RESULTS = 'CLEAR_SEARCH_RESULTS';
export function clearSearchResults() {
    return {
        type: CLEAR_SEARCH_RESULTS
    }
}

export const RECEIVE_SEARCH_RESULTS_SUCCESS = 'RECEIVE_SEARCH_RESULTS_SUCCESS';
function receiveSearchResultsSuccess(timestamp, results) {
    return {
        type: RECEIVE_SEARCH_RESULTS_SUCCESS,
        timestamp,
        results
    }
}
