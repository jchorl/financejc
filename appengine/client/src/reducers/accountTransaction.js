import Immutable from 'immutable';
import {
	PUT_TRANSACTION,
	RECEIVE_TRANSACTIONS,
	RECEIVE_ACCOUNTS
} from '../actions'

export default (state = Immutable.Map(), action) => {
	switch (action.type) {
		case RECEIVE_TRANSACTIONS:
			// make the dates into actual dates
			let transactions = Immutable.fromJS(action.transactions);
			transactions = transactions.map(t => t.setIn(['date'], new Date(t.get('date'))));
			state = state.setIn([action.accountId, "transactions"], state.get(action.accountId).get("transactions").concat(transactions));

			let re = new RegExp("<(.*)>; rel=\"next\"");
			let result = re.exec(action.link)
			if (result) {
				state = state.setIn([action.accountId, "next"], result[1]);
			}

			return state;
		case PUT_TRANSACTION:
			let transaction = Immutable.fromJS(action.transaction);
			transaction = transaction.set('date', new Date(transaction.get('date')));

			// TODO update properly
			return state.get(transaction.accountId).unshift(transaction);
		case RECEIVE_ACCOUNTS:
			return state.withMutations(map => {
				for (let account of action.accounts) {
					map = map.set(account.id, Immutable.Map({
						"next": null,
						"transactions": Immutable.List()
					}));
				}
			});
		default:
			return state
	}
}
