import Immutable from 'immutable';
import {
	RECEIVE_ACCOUNTS
} from '../actions'

export default (state = Immutable.Map(), action) => {
	switch (action.type) {
		case RECEIVE_ACCOUNTS:
			return Immutable.Map().withMutations(map => {
				for (let account of action.accounts) {
					map.set(account.id, Immutable.fromJS(account));
				}
			});
		default:
			return state;
	}
}