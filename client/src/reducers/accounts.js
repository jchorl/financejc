import Immutable from 'immutable';
import {
  ADD_ACCOUNT,
  RECEIVE_ACCOUNTS,
  UPDATE_ACCOUNT_VALUE,
  LOGOUT
} from '../actions';

export default (state = Immutable.Map({
    fetched: false,
    accounts: Immutable.Map()
}), action) => {
    switch (action.type) {
    case ADD_ACCOUNT:
        return state.setIn(['accounts', action.account.id], Immutable.fromJS(action.account));

    case RECEIVE_ACCOUNTS:
        return state.set('fetched', true).set('accounts', Immutable.Map().withMutations(map => {
            for (let account of action.accounts) {
                map.set(account.id, Immutable.fromJS(account));
            }
        }));

    case UPDATE_ACCOUNT_VALUE:
        return state.updateIn(['accounts', action.accountId, 'futureValue'], val => val + action.delta);

    case LOGOUT:
        return Immutable.Map({
            fetched: false,
            accounts: Immutable.Map()
        });

    default:
        return state;
    }
};
