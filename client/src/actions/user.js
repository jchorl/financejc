import { handleErrors } from './util';

export function fetchUser() {
    return function(dispatch) {
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        return fetch('/api/user', {
            headers,
            credentials: 'include'
        })
        .then(handleErrors)
            .then(resp => resp.json())
            .then(user => {
                dispatch(receiveUserSuccess(user));
                return user;
            })
        .catch(e => {
            dispatch(receiveUserFailure());
            console.error(e)
        });
    }
}

export const RECEIVE_USER_SUCCESS = 'RECEIVE_USER_SUCCESS';
function receiveUserSuccess(user) {
    return {
        type: RECEIVE_USER_SUCCESS,
        user
    }
}

export const RECEIVE_USER_FAILURE = 'RECEIVE_USER_FAILURE';
function receiveUserFailure() {
    return {
        type: RECEIVE_USER_FAILURE
    }
}

export function googleLogin(googleUser) {
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
        .then(handleErrors)
            .then(response => response.json())
            .then(user => {
                dispatch(receiveUserSuccess(user));
                return user;
            })
        .catch(e => {
            dispatch(receiveUserFailure());
            console.error(e)
        });
    };
}
