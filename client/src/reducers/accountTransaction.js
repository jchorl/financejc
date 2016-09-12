import Immutable from 'immutable';
import {
	PUT_TRANSACTION,
	RECEIVE_TRANSACTIONS,
	RECEIVE_ACCOUNTS
} from '../actions'

export default (state = Immutable.Map(), action) => {
	switch (action.type) {
		case RECEIVE_TRANSACTIONS:
			let transactions = Immutable.fromJS(action.transactions).reduce(
				function(result, item) {
					item = item.set('date', new Date(item.get('date')))
					return result.set(item.get('id'), item);
				}, 
				Immutable.Map()
			);

			state = state.setIn([action.accountId, 'transactions'], state.get(action.accountId).get('transactions').merge(transactions).sortBy(t => -t.get('date')));

			let re = new RegExp('<(.*)>; rel="next"');
			let result = re.exec(action.link)
			if (result) {
				state = state.setIn([action.accountId, 'next'], result[1]);
			}

			return state;

		case PUT_TRANSACTION:
			let transaction = Immutable.fromJS(action.transaction);
			transaction = transaction.set('date', new Date(transaction.get('date')));
			let updated = state.get(transaction.get('account')).get('transactions').set(transaction.get('id'), transaction).sortBy(t => -t.get('date'));
			return state.setIn([transaction.get('account'), 'transactions'], updated);

		case RECEIVE_ACCOUNTS:
			return state.withMutations(map => {
				for (let account of action.accounts) {
					map = map.set(account.id, Immutable.Map({
						'next': null,
						'transactions': Immutable.Map()
					}));
				}
			});

		default:
			return state
	}
}
