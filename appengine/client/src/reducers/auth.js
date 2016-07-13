import {
	CHECK_AUTH,
	RECEIVE_AUTH,
	REQUEST_LOGIN
} from '../actions'

export default (state = {
	isFetching: false,
	authd: false
}, action) => {
	switch (action.type) {
		case CHECK_AUTH:
			return Object.assign({}, state, {
				isFetching: true
			});
		case REQUEST_LOGIN:
			return Object.assign({}, state, {
				isFetching: true
			});
		case RECEIVE_AUTH:
			return Object.assign({}, state, {
				isFetching: false,
				authd: action.authd
			});
		default:
			return state
	}
}

