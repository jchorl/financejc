import { List, Map, fromJS } from 'immutable';
import { CLEAR_SEARCH_RESULTS, RECEIVE_SEARCH_RESULTS_SUCCESS } from '../actions/search';
import { EDIT_TRANSACTION_SUCCESS, DELETE_TRANSACTION_SUCCESS } from '../actions/accountData';
import { dateStringToDate } from './util';

export default function(state = Map({
    lastSearchTimestamp: new Date(),
    items: List()
}), action) {
    switch (action.type) {
        case RECEIVE_SEARCH_RESULTS_SUCCESS:
            // ensure most recent query is always kept
            if (action.timestamp > state.get('lastSearchTimestamp')) {
                state = state.set('lastSearchTimestamp', action.timestamp);
                state = state.set('items', fromJS(action.results).map(t => t.update('date', d => dateStringToDate(d))))
            }
            return state;

        case CLEAR_SEARCH_RESULTS:
            return Map({
                lastSearchTimestamp: new Date(),
                items: List()
            })

        case EDIT_TRANSACTION_SUCCESS: {
            let index = state.get('items').findIndex(t => t.get('id') === action.transaction.id);
            if (index !== -1) {
                let parsed = fromJS(action.transaction).update('date', d => dateStringToDate(d));
                state = state.setIn(['items', index], parsed);
            }
            return state;
        }

        case DELETE_TRANSACTION_SUCCESS:
            return state.update('items', items => items.filter(t => t.get('id') !== action.transaction.id));

        default:
            return state;
    }
}
