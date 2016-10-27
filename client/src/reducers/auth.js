import Immutable from 'immutable';
import {
  CHECK_AUTH,
  RECEIVE_AUTH,
  REQUEST_LOGIN,
  LOGOUT
} from '../actions';

export default (state = Immutable.Map({
  fetched: false,
  authd: false
}), action) => {
  switch (action.type) {
    case CHECK_AUTH:
      return state.set('fetched', false);

    case REQUEST_LOGIN:
      return state.set('fetched', false);

    case RECEIVE_AUTH:
      return Immutable.Map({
        fetched: true,
        authd: action.authd
      });

    case LOGOUT:
      return state.set('authd', false);

    default:
      return state
  }
}

