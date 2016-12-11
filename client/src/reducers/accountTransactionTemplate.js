import Immutable from 'immutable';
import {
  ADD_ACCOUNT,
  PUT_TRANSACTION_TEMPLATE,
  RECEIVE_TRANSACTION_TEMPLATES,
  RECEIVE_ACCOUNTS,
  LOGOUT
} from '../actions';

export default (state = Immutable.Map(), action) => {
    switch (action.type) {
    case RECEIVE_TRANSACTION_TEMPLATES: {
        let transactions = Immutable.fromJS(action.transactionTemplates).reduce((result, item) => result.set(item.get('id'), item), Immutable.Map());

        state = state.setIn([action.accountId, 'fetched'], true);
        return state.updateIn([action.accountId, 'transactionTemplates'], ts => ts.merge(transactions).sortBy(t => t.get('name')));
    }

    case PUT_TRANSACTION_TEMPLATE: {
        let transaction = Immutable.fromJS(action.transactionTemplate);
        return state.updateIn([transaction.get('accountId'), 'transactionTemplates'], transactions => transactions.set(transaction.get('id'), transaction).sortBy(t => t.getIn('name')));
    }

    case ADD_ACCOUNT:
        return state.set(action.account.id, Immutable.Map({
            'fetched': false,
            'transactionTemplates': Immutable.Map()
        }));

    case RECEIVE_ACCOUNTS:
        return state.withMutations(map => {
            for (let account of action.accounts) {
                map = map.set(account.id, Immutable.Map({
                    'fetched': false,
                    'transactionTemplates': Immutable.Map()
                }));
            }
        });

    case LOGOUT:
        return Immutable.Map();

    default:
        return state;
    }
};
