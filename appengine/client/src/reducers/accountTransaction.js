import Immutable from 'immutable';
import accounts from './accounts';
import transactions from './transactions';
import {
	PUT_TRANSACTION
} from '../actions'

function insert(transaction, transactionList, transactionState) {
  return transactionList.insert(locationOf(transaction, transactionList, transactionState, dateCompare) + 1, transaction.get('id'));
}

function locationOf(transaction, transactionList, transactionState, comparer, start, end) {
    if (transactionList.size === 0)
        return -1;

    start = start || 0;
    end = end || transactionList.size;
    let pivot = (start + end) >> 1;  // should be faster than the above calculation

    let c = comparer(transactionState.get(transaction.get('id')), transactionState.get(transactionList.get(pivot)));
    if (end - start <= 1) return c == -1 ? pivot - 1 : pivot;

    switch (c) {
        case -1: return locationOf(transaction, transactionList, transactionState, comparer, start, pivot);
        case 0: return pivot;
        case 1: return locationOf(transaction, transactionList, transactionState, comparer, pivot, end);
    };
};

function dateCompare(a, b) {
	if (a.get('date') > b.get('date')) return -1;
	if (a.get('date') < b.get('date')) return 1;
	return 0;
};

// reorder account transactions
function reorderAccountTransactions(transaction, accountState, transactionState) {
	let account = accountState.get(transaction.get('accountId')).asMutable();
	// need to make sure the transaction gets inserted in the right place
	// first remove
	let accountTransactions = account.get('transactions');
	let index = accountTransactions.findIndex(id => id === transaction.get('id'));
	if (index > -1) {
		accountTransactions = accountTransactions.delete(index);
	}

	accountTransactions = insert(transaction, accountTransactions, transactionState);
	account.set('transactions', accountTransactions);
	return accountState.set(account.get('id'), account.asImmutable());
}

export default (state = Immutable.Map(), action) => {
	let accountData = accounts(state.get('account'), action);
	const transactionData = transactions(state.get('transaction'), action);
	if (action.type === PUT_TRANSACTION) {
		// console.log('before: ');
		// let z = accountData["aghkZXZ-Tm9uZXIlCxIEVXNlchiAgICAgICACgwLEgdBY2NvdW50GICAgICAgJsIDA"].transactions[0];
		// let f = accountData["aghkZXZ-Tm9uZXIlCxIEVXNlchiAgICAgICACgwLEgdBY2NvdW50GICAgICAgJsIDA"].transactions[1];
		// console.log(z);
		// console.log(transactionData[z]);
		// console.log(f);
		// console.log(transactionData[f]);
		accountData = reorderAccountTransactions(transactionData.get(action.transaction.id), accountData, transactionData);
		// console.log('after: ');
		// z = accountData["aghkZXZ-Tm9uZXIlCxIEVXNlchiAgICAgICACgwLEgdBY2NvdW50GICAgICAgJsIDA"].transactions[0];
		// f = accountData["aghkZXZ-Tm9uZXIlCxIEVXNlchiAgICAgICACgwLEgdBY2NvdW50GICAgICAgJsIDA"].transactions[1];
		// console.log(z);
		// console.log(transactionData[z]);
		// console.log(f);
		// console.log(transactionData[f]);
	}
	return Immutable.fromJS({
		account: accountData,
		transaction: transactionData
	});
}
