import { List, Map, fromJS } from 'immutable';
import { RECEIVE_SUGGESTIONS_SUCCESS } from '../actions/autocomplete';

export default function(state = Map({
    field: '',
    term: '',
    items: List()
}), action) {
    switch (action.type) {
        case RECEIVE_SUGGESTIONS_SUCCESS:
            return Map({
                field: action.field,
                term: action.term,
                items: fromJS(action.suggestions)
            });

        default:
            return state;
    }
}
