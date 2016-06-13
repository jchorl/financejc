import {
	UPDATE_TRANSACTION
} from '../actions'

export default (state = {
	transaction: {}
}, action) => {
	switch (action.type) {
		case UPDATE_TRANSACTION:
			return Object.assign({}, state, {
				transaction: action.transaction
			});
		default:
			return state
	}
}

