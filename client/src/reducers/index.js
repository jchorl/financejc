import { combineReducers } from 'redux';
import auth from './auth';
import accounts from './accounts';
import accountTransaction from './accountTransaction';
import currencies from './currencies';

export default combineReducers({
  auth,
  accounts,
  accountTransaction,
  currencies
});
