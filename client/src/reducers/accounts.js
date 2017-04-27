import Immutable from 'immutable';
import {
    ADD_ACCOUNT,
    RECEIVE_ACCOUNTS,
    SELECT_ACCOUNT,
    UPDATE_ACCOUNT_VALUE,
    LOGOUT
} from '../actions';

export default (state = Immutable.Map({
    fetched: false,
    accounts: Immutable.OrderedMap(),
    selected: -1
}), action) => {
    switch (action.type) {
    case ADD_ACCOUNT:
        return state.setIn(['accounts', action.account.id], Immutable.fromJS(action.account));

    case SELECT_ACCOUNT:
        return state.set('selected', action.id);

    case RECEIVE_ACCOUNTS: {
        let sortedAccounts = Immutable.OrderedMap(
                Immutable.Seq(action.accounts)
                .sort((c1, c2) => c1.name > c2.name)
                .map(account => [account.id, Immutable.fromJS(account)])
                );
        return state
                .set('fetched', true)
                .set('accounts', sortedAccounts);
    }

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
