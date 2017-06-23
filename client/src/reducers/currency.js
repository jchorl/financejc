import { Map, fromJS } from 'immutable';
import { RECEIVE_CURRENCIES_SUCCESS } from '../actions/currency';

export default function(state = Map({
    fetched: false,
    items: Map()
}), action) {
    switch (action.type) {
        case RECEIVE_CURRENCIES_SUCCESS:
            return Map({
                fetched: true,
                items: fromJS(action.currencies)
            });

        default:
            return state;
    }
}
