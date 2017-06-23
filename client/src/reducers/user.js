import { Map } from 'immutable';
import { RECEIVE_USER_SUCCESS, RECEIVE_USER_FAILURE } from '../actions/user';
import { UNAUTHD_ID } from '../constants';

export default function(state = Map({
    fetched: false,
    email: '',
    id: UNAUTHD_ID
}), action) {
    switch (action.type) {
        case RECEIVE_USER_SUCCESS:
            const {
                email,
                id
            } = action.user;

            return Map({
                fetched: true,
                email,
                id
            });

        case RECEIVE_USER_FAILURE:
            return Map({
                fetched: true,
                email: '',
                id: UNAUTHD_ID
            });

        default:
            return state;
    }
}
