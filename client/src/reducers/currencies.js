import { fromJS, Map, Seq } from 'immutable';
import {
  RECEIVE_CURRENCIES
} from '../actions';

export default (state = Map({
    fetched: false,
    currencies: Map()
}), action) => {
    switch (action.type) {
    case RECEIVE_CURRENCIES: {
        let sortedCurrencies = Seq(action.currencies).map(v => fromJS(v)).toOrderedMap();
        return state.set('fetched', true).set('currencies', sortedCurrencies);
    }
    default:
        return state;
    }
};
