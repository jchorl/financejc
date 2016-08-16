import Immutable from 'immutable';
import {
	RECEIVE_TRANSACTIONS,
	PUT_TRANSACTION
} from '../actions'

export default (state = Immutable.Map(), action) => {
	switch (action.type) {
		case RECEIVE_TRANSACTIONS:
			// make the dates into actual dates
			let transactions = Immutable.fromJS(action.transactions);
			transactions = transactions.map(t => t.setIn(['date'], new Date(t.get('date'))));
			return transactions;
		case PUT_TRANSACTION:
			let transaction = Immutable.fromJS(action.transaction);
			transaction = transaction.set('date', new Date(transaction.get('date')));
			return state.set(transaction.get('id'), transaction);
		default:
			return state
	}
}

