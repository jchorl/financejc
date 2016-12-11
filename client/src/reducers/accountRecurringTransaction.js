import Immutable from 'immutable';
import {
  ADD_ACCOUNT,
  PUT_RECURRING_TRANSACTION,
    DELETE_RECURRING_TRANSACTION,
  RECEIVE_RECURRING_TRANSACTIONS,
  RECEIVE_ACCOUNTS,
  LOGOUT
} from '../actions';

export default (state = Immutable.Map(), action) => {
    switch (action.type) {
    case RECEIVE_RECURRING_TRANSACTIONS: {
        let transactions = Immutable.fromJS(action.recurringTransactions).reduce(
        function(result, item) {
            item = item.setIn(['transaction', 'date'], new Date(item.get('transaction').get('date')));
            return result.set(item.get('id'), item);
        },
        Immutable.Map()
      );

        state = state.setIn([action.accountId, 'fetched'], true);
        return state.updateIn([action.accountId, 'recurringTransactions'], ts => ts.merge(transactions).sortBy(t => -t.getIn(['transaction', 'date'])));
    }

    case PUT_RECURRING_TRANSACTION: {
        let transaction = Immutable.fromJS(action.recurringTransaction);
        transaction = transaction.setIn(['transaction', 'date'], new Date(transaction.getIn(['transaction', 'date'])));
        return state.updateIn([transaction.getIn(['transaction', 'accountId']), 'recurringTransactions'], transactions => transactions.set(transaction.get('id'), transaction).sortBy(t => -t.getIn(['transaction', 'date'])));
    }

    case DELETE_RECURRING_TRANSACTION: {
        return state.updateIn([action.accountId, 'recurringTransactions'], transactions => transactions.delete(action.id));
    }

    case ADD_ACCOUNT:
        return state.set(action.account.id, Immutable.Map({
            'fetched': false,
            'recurringTransactions': Immutable.Map()
        }));

    case RECEIVE_ACCOUNTS:
        return state.withMutations(map => {
            for (let account of action.accounts) {
                map = map.set(account.id, Immutable.Map({
                    'fetched': false,
                    'recurringTransactions': Immutable.Map()
                }));
            }
        });

    case LOGOUT:
        return Immutable.Map();

    default:
        return state;
    }
};
