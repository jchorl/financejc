import Immutable from 'immutable';
import {
    FETCHING_USER,
    RECEIVE_USER,
    LOGOUT
} from '../actions';

export default (state = Immutable.Map({
    fetched: false,
    user: Immutable.Map(),
    authd: false
}), action) => {
    switch (action.type) {
    case FETCHING_USER:
        return state.set('fetched', false);

    case RECEIVE_USER:
        if (action.user) {
            return Immutable.Map({
                fetched: true,
                user: Immutable.fromJS(action.user),
                authd: true
            });
        }

        return Immutable.Map({
            fetched: true,
            user: Immutable.Map(),
            authd: false
        });

    case LOGOUT:
        return Immutable.Map({
            fetched: true,
            user: Immutable.Map(),
            authd: false
        });

    default:
        return state;
    }
};

