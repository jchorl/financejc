import { fromJS, List } from 'immutable';
import { RECEIVE_SUMMARY_SUCCESS } from '../actions/summary';
import { dateStringToDate } from './util';


export default function(state = List(), action) {
    switch (action.type) {
        case RECEIVE_SUMMARY_SUCCESS:
            return fromJS(action.transactions).map(t => t.update('date', d => dateStringToDate(d)));

        default:
            return state;
    }
}
