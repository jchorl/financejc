import {
	RECEIVE_ACCOUNTS
} from '../actions'

export default (state = {}, action) => {
	switch (action.type) {
		case RECEIVE_ACCOUNTS:
			return action.accounts
		default:
			return state
	}
}
