import Immutable from 'immutable';
import {
    ADD_ACCOUNT,
    DELETE_TRANSACTION,
    PUT_TRANSACTION,
    RECEIVE_TRANSACTIONS,
    RECEIVE_ACCOUNTS,
    LOGOUT
} from '../actions';

export default (state = Immutable.Map(), action) => {
    switch (action.type) {
    case RECEIVE_TRANSACTIONS: {
        let transactions = Immutable.fromJS(action.transactions).reduce(
                function(result, item) {
                    item = item.set('date', new Date(item.get('date')));
                    return result.set(item.get('id'), item);
                },
                Immutable.Map()
            );

        state = state.updateIn([action.accountId, 'transactions'], ts => ts.merge(transactions).sortBy(t => -t.get('date')));

        let re = new RegExp('<(.*)>; rel="next"');
        let result = re.exec(action.link);
        if (result) {
            state = state.setIn([action.accountId, 'next'], result[1]);
        }

        return state;
    }

    case PUT_TRANSACTION: {
        let transaction = Immutable.fromJS(action.transaction);
        transaction = transaction.set('date', new Date(transaction.get('date')));
        return state.updateIn([transaction.get('accountId'), 'transactions'], transactions => transactions.set(transaction.get('id'), transaction).sortBy(t => -t.get('date')));
    }

    case DELETE_TRANSACTION: {
        return state.updateIn([action.accountId, 'transactions'], transactions => transactions.delete(action.id));
    }

    case ADD_ACCOUNT:
        return state.set(action.account.id, Immutable.Map({
            'next': null,
            'transactions': Immutable.Map()
        }));

    case RECEIVE_ACCOUNTS:
        return state.withMutations(map => {
            for (let account of action.accounts) {
                map = map.set(account.id, Immutable.Map({
                    'next': null,
                    'transactions': Immutable.Map()
                }));
            }
        });

    case LOGOUT:
        return Immutable.Map();

    default:
        return state;
    }
};
