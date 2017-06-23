import { handleErrors } from './util';

export function fetchSuggestions(accountId, field, term) {
    return function(dispatch) {
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        return fetch(`/api/account/${accountId}/transactions/query?field=${field}&value=${term}`, {
            headers,
            credentials: 'include'
        })
        .then(handleErrors)
            .then(resp => resp.json())
            .then(suggestions => {
                dispatch(receiveSuggestionsSuccess(field, term, suggestions));
                return suggestions;
            })
        .catch(e => console.error(e));
    }
}

export const RECEIVE_SUGGESTIONS_SUCCESS = 'RECEIVE_SUGGESTIONS_SUCCESS';
function receiveSuggestionsSuccess(field, term, suggestions) {
    return {
        type: RECEIVE_SUGGESTIONS_SUCCESS,
        field,
        term,
        suggestions
    }
}
