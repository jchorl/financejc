import { combineReducers } from 'redux';
import {reducer as form} from 'redux-form';
import accounts from './accounts';
import auth from './auth';
import fetching from './fetching';
import transactions from './transactions';

export default combineReducers({
	fetching,
	accounts,
	auth,
	transactions,
	form
});