import Immutable from 'immutable';
import {
  RECEIVE_ACCOUNTS,
  UPDATE_ACCOUNT_VALUE,
  LOGOUT
} from '../actions';

export default (state = Immutable.Map(), action) => {
  switch (action.type) {
    case RECEIVE_ACCOUNTS:
      return Immutable.Map().withMutations(map => {
        for (let account of action.accounts) {
          map.set(account.id, Immutable.fromJS(account));
        }
      });

    case UPDATE_ACCOUNT_VALUE:
      return state.updateIn([action.accountId, 'futureValue'], val => val + action.delta);

    case LOGOUT:
      return Immutable.Map()

    default:
      return state;
  }
}
