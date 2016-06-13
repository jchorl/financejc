import { combineReducers } from 'redux';
import {reducer as form} from 'redux-form';
import accounts from './accounts';
import auth from './auth';
import transactions from './transactions';

export default combineReducers({
	accounts,
	auth,
	transactions,
	form
});