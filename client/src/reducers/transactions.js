import {
	RECEIVE_TRANSACTIONS,
	UPDATE_TRANSACTION
} from '../actions'

export default (state = {}, action) => {
	switch (action.type) {
		case RECEIVE_TRANSACTIONS:
			return action.transactions;
		case UPDATE_TRANSACTION:
			updated = Object.assign({}, state);
			updated[action.transaction.id] = action.transaction;
			return updated;
		default:
			return state
	}
}

