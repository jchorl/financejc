import Immutable from 'immutable';
import {
    ADD_ACCOUNT,
    PUT_TEMPLATE,
    DELETE_TEMPLATE,
    RECEIVE_TEMPLATES,
    RECEIVE_ACCOUNTS,
    LOGOUT
} from '../actions';

export default (state = Immutable.Map(), action) => {
    switch (action.type) {
    case RECEIVE_TEMPLATES: {
        let transactions = Immutable.fromJS(action.templates).reduce((result, item) => result.set(item.get('id'), item), Immutable.Map());

        state = state.setIn([action.accountId, 'fetched'], true);
        return state.updateIn([action.accountId, 'templates'], ts => ts.merge(transactions).sortBy(t => t.get('name')));
    }

    case PUT_TEMPLATE: {
        let transaction = Immutable.fromJS(action.template);
        return state.updateIn([transaction.get('accountId'), 'templates'], transactions => transactions.set(transaction.get('id'), transaction).sortBy(t => t.getIn('name')));
    }

    case DELETE_TEMPLATE: {
        return state.updateIn([action.accountId, 'templates'], transactions => transactions.delete(action.id));
    }

    case ADD_ACCOUNT:
        return state.set(action.account.id, Immutable.Map({
            'fetched': false,
            'templates': Immutable.Map()
        }));

    case RECEIVE_ACCOUNTS:
        return state.withMutations(map => {
            for (let account of action.accounts) {
                map = map.set(account.id, Immutable.Map({
                    'fetched': false,
                    'templates': Immutable.Map()
                }));
            }
        });

    case LOGOUT:
        return Immutable.Map();

    default:
        return state;
    }
};
