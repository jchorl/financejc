import {
	RECEIVE_TRANSACTIONS,
	PUT_TRANSACTION
} from '../actions'

export default (state = {}, action) => {
	switch (action.type) {
		case RECEIVE_TRANSACTIONS:
			return action.transactions;
		case PUT_TRANSACTION:
			let updated = Object.assign({}, state);
			updated[action.transaction.id] = action.transaction;
			return updated;
		default:
			return state
	}
}

