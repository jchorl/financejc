import Immutable from 'immutable';
import {
	PUT_TRANSACTION,
	RECEIVE_TRANSACTIONS,
	RECEIVE_ACCOUNTS
} from '../actions'

let transactionCompare = function (a, b) {
    if (a.get('date') > b.get('date')) return -1;
    if (a.get('date') < b.get('date')) return 1;
    return 0;
};

function insert(element, array) {
  return array.insert(locationOf(element, array, transactionCompare) + 1, element);
}

function locationOf(element, array, comparer, start, end) {
    if (array.size === 0)
        return -1;

    start = start || 0;
    end = end || array.size;
    let pivot = (start + end) >> 1;  // should be faster than the above calculation

    let c = comparer(element, array.get(pivot));
    if (end - start <= 1) return c == -1 ? pivot - 1 : pivot;

    switch (c) {
        case -1: return locationOf(element, array, comparer, start, pivot);
        case 0: return pivot;
        case 1: return locationOf(element, array, comparer, pivot, end);
    };
};

export default (state = Immutable.Map(), action) => {
	switch (action.type) {
		case RECEIVE_TRANSACTIONS:
			// make the dates into actual dates
			let transactions = Immutable.fromJS(action.transactions);
			transactions = transactions.map(t => t.setIn(['date'], new Date(t.get('date'))));
			state = state.setIn([action.accountId, 'transactions'], state.get(action.accountId).get('transactions').concat(transactions));

			let re = new RegExp('<(.*)>; rel="next"');
			let result = re.exec(action.link)
			if (result) {
				state = state.setIn([action.accountId, 'next'], result[1]);
			}

			return state;
		case PUT_TRANSACTION:
			let transaction = Immutable.fromJS(action.transaction);
			transaction = transaction.set('date', new Date(transaction.get('date')));
			let transactionId = transaction.get('id');
			let updated = state.get(transaction.get('accountId')).get('transactions').asMutable().filter(t => t.get('id') != transactionId);
			updated = insert(transaction, updated);
			return state.setIn([transaction.get('accountId'), 'transactions'], updated.asImmutable());
		case RECEIVE_ACCOUNTS:
			return state.withMutations(map => {
				for (let account of action.accounts) {
					map = map.set(account.id, Immutable.Map({
						'next': null,
						'transactions': Immutable.List()
					}));
				}
			});
		default:
			return state
	}
}
