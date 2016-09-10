import { combineReducers } from 'redux';
import {reducer as form} from 'redux-form';
import auth from './auth';
import fetching from './fetching';
import account from './accounts';
import accountTransaction from './accountTransaction.js';

export default combineReducers({
	fetching,
	auth,
	account,
	accountTransaction,
	form
});