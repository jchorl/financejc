import Immutable from 'immutable';
import {
  RECEIVE_CURRENCIES
} from '../actions';

export default (state = Immutable.Map({
    fetched: false,
    currencies: Immutable.Map()
}), action) => {
    switch (action.type) {
    case RECEIVE_CURRENCIES: {
        let sortedCurrencies = Immutable.Seq(action.currencies)
            .sort((c1, c2) => c1.code > c2.code)
            .map(v => Immutable.fromJS(v))
            .toOrderedMap();
        return state.set('fetched', true).set('currencies', sortedCurrencies);
    }
    default:
        return state;
    }
};
