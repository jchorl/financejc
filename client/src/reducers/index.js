import { combineReducers } from 'redux';
import auth from './auth';
import accounts from './accounts';
import accountTransaction from './accountTransaction';
import accountRecurringTransaction from './accountRecurringTransaction';
import accountTransactionTemplate from './accountTransactionTemplate';
import currencies from './currencies';

export default combineReducers({
    auth,
    accounts,
    accountTransaction,
    accountRecurringTransaction,
    accountTransactionTemplate,
    currencies
});
