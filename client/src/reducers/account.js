import { Map, Seq, fromJS } from 'immutable';
import { CREATE_ACCOUNT_SUCCESS, RECEIVE_ACCOUNTS_SUCCESS } from '../actions/account';
import { EDIT_TRANSACTION_SUCCESS, DELETE_TRANSACTION_SUCCESS } from '../actions/accountData';

export default function(state = Map({
    fetched: false,
    items: Map()
}), action) {
    switch (action.type) {
        case RECEIVE_ACCOUNTS_SUCCESS:
            return Map({
                fetched: true,
                items: Seq(action.accounts).reduce(
                        (m, acc) => m.set(acc.id, fromJS(acc)),
                        Map()
                        )
            });

        case CREATE_ACCOUNT_SUCCESS:
            return state.setIn(['items', action.account.id], fromJS(action.account));

        case EDIT_TRANSACTION_SUCCESS:
            return state.updateIn(['items', action.transaction.accountId, 'futureValue'], fv => fv + action.amountDifference);

        case DELETE_TRANSACTION_SUCCESS:
            return state.updateIn(['items', action.transaction.accountId, 'futureValue'], fv => fv - action.transaction.amount);

        default:
            return state;
    }
}
