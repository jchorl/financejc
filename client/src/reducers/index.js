import { combineReducers } from 'redux';
import account from './account';
import accountData from './accountData';
import autocomplete from './autocomplete';
import currency from './currency';
import search from './search';
import summary from './summary';
import user from './user';

export default combineReducers({
    account,
    accountData,
    autocomplete,
    currency,
    search,
    summary,
    user
});
