import {
  REQUEST_ACCOUNTS,
  RECEIVE_ACCOUNTS,
  LOGOUT
} from '../actions';

export default (state = false, action) => {
  switch (action.type) {
    case REQUEST_ACCOUNTS:
      return true;

    case RECEIVE_ACCOUNTS:
      return false;

    default:
      return state
  }
}
