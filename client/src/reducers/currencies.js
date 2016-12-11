import Immutable from 'immutable';
import {
  RECEIVE_CURRENCIES
} from '../actions';

export default (state = Immutable.Map({
    fetched: false,
    currencies: Immutable.Map()
}), action) => {
    switch (action.type) {
    case RECEIVE_CURRENCIES:
        return state.set('fetched', true).set('currencies', Immutable.fromJS(action.currencies));
    default:
        return state;
    }
};
