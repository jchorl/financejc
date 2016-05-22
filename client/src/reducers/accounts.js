import {
	REQUEST_ACCOUNTS,
	RECEIVE_ACCOUNTS
} from '../actions'

export default (state = {
	isFetching: false,
	items: []
}, action) => {
	switch (action.type) {
		case REQUEST_ACCOUNTS:
			return Object.assign({}, state, {
				isFetching: true
			});
		case RECEIVE_ACCOUNTS:
			return Object.assign({}, state, {
				isFetching: false,
				items: action.accounts
			});
		default:
			return state
	}
}
