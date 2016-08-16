import Immutable from 'immutable';
import {
	RECEIVE_ACCOUNTS
} from '../actions'

export default (state = Immutable.Map(), action) => {
	switch (action.type) {
		case RECEIVE_ACCOUNTS:
			return Immutable.fromJS(action.accounts);
		default:
			return state;
	}
}
