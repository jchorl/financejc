import { combineReducers } from 'redux';
import {reducer as form} from 'redux-form';
import auth from './auth';
import fetching from './fetching';
import accountTransaction from './accountTransaction.js';

export default combineReducers({
	fetching,
	auth,
	accountTransaction,
	form
});