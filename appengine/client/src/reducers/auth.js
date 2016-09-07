import {
	CHECK_AUTH,
	RECEIVE_AUTH,
	REQUEST_LOGIN
} from '../actions'

export default (state = {
	fetched: false,
	authd: false
}, action) => {
	switch (action.type) {
		case CHECK_AUTH:
			return Object.assign({}, state, {
				fetched: false
			});
		case REQUEST_LOGIN:
			return Object.assign({}, state, {
				fetched: false
			});
		case RECEIVE_AUTH:
			return Object.assign({}, state, {
				fetched: true,
				authd: action.authd
			});
		default:
			return state
	}
}

