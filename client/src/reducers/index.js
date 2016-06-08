import { combineReducers } from 'redux';
import {reducer as form} from 'redux-form';
import accounts from './accounts';
import auth from './auth';

export default combineReducers({
	accounts,
	auth,
	form
});