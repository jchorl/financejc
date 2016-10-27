import { combineReducers } from 'redux';
import {reducer as form} from 'redux-form';
import auth from './auth';
import accounts from './accounts';
import accountTransaction from './accountTransaction';
import currencies from './currencies';

export default combineReducers({
  auth,
  accounts,
  accountTransaction,
  form,
  currencies
});
