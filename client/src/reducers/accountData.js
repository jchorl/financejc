import { Map, List, Seq, fromJS } from 'immutable';
import { CREATE_ACCOUNT_SUCCESS, RECEIVE_ACCOUNTS_SUCCESS } from '../actions/account';
import {
    EDIT_TRANSACTION_SUCCESS,
    DELETE_TRANSACTION_SUCCESS,
    FETCH_TRANSACTIONS_START,
    RECEIVE_TRANSACTIONS_SUCCESS,
    EDIT_TEMPLATE_SUCCESS,
    FETCH_TEMPLATES_START,
    RECEIVE_TEMPLATES_SUCCESS,
    EDIT_RECURRING_TRANSACTION_SUCCESS,
    DELETE_RECURRING_TRANSACTION_SUCCESS,
    FETCH_RECURRING_TRANSACTIONS_START,
    RECEIVE_RECURRING_TRANSACTIONS_SUCCESS
} from '../actions/accountData';
import { dateStringToDate } from './util';

function emptyAccountData(account) {
    return Map({
        transactions: Map({
            fetched: false,
            fetching: false,
            futureValue: account.futureValue,
            data: Map(),
            calculatedValues: List()
        }),
        templates: Map({
            fetched: false,
            fetching: false,
            items: List()
        }),
        recurringTransactions: Map({
            fetched: false,
            fetching: false,
            items: List()
        })
    });
}

export default function(state = Map(), action) {
    switch (action.type) {
        case RECEIVE_ACCOUNTS_SUCCESS:
            return Seq(action.accounts)
                .reduce(
                        (m, acc) => m.set(acc.id, emptyAccountData(acc)),
                        Map()
                        );

        case CREATE_ACCOUNT_SUCCESS:
            return state.set(action.account.id, emptyAccountData(action.account));

        case FETCH_TRANSACTIONS_START:
            return state.setIn([action.accountId, 'transactions', 'fetching'], true);

        case RECEIVE_TRANSACTIONS_SUCCESS: {
            let data = Seq(action.transactions)
                .reduce((m, t) => {
                    let parsed = fromJS(t);
                    parsed = parsed.update('date', d => dateStringToDate(d));
                    return m.set(t.id, parsed);
                }, Map());
            data = state.getIn([action.accountId, 'transactions', 'data']).merge(data);
            let futureValue = state.getIn([action.accountId, 'transactions', 'futureValue']);
            let calculatedValues = calculateValues(data, futureValue);
            return state.setIn([action.accountId, 'transactions'], Map({
                data,
                calculatedValues,
                futureValue,
                fetched: true,
                fetching: false
            }));
        }

        case EDIT_TRANSACTION_SUCCESS: {
            let parsed = fromJS(action.transaction);
            parsed = parsed.update('date', d => dateStringToDate(d));

            let accountId = parsed.get('accountId');
            state = state.setIn([accountId, 'transactions', 'data', parsed.get('id')], parsed);
            state = state.updateIn([accountId, 'transactions', 'futureValue'], fv => fv + action.amountDifference);
            state = state.setIn([accountId, 'transactions', 'calculatedValues'], calculateValues(
                        state.getIn([accountId, 'transactions', 'data']),
                        state.getIn([accountId, 'transactions', 'futureValue']),
                        ));

            return state;
        }

        case DELETE_TRANSACTION_SUCCESS: {
            let transaction = action.transaction;
            let accountId = transaction.accountId;
            state = state.deleteIn([accountId, 'transactions', 'data', transaction.id]);
            state = state.updateIn([accountId, 'transactions', 'futureValue'], fv => fv - transaction.amount);
            state = state.setIn([accountId, 'transactions', 'calculatedValues'], calculateValues(
                        state.getIn([accountId, 'transactions', 'data']),
                        state.getIn([accountId, 'transactions', 'futureValue']),
                        ));

            return state;
        }

        case FETCH_TEMPLATES_START:
            return state.setIn([action.accountId, 'templates', 'fetching'], true);

        case RECEIVE_TEMPLATES_SUCCESS:
            return state.setIn([action.accountId, 'templates'], Map({
                items: fromJS(action.templates),
                fetched: true,
                fetching: false
            }));

        case EDIT_TEMPLATE_SUCCESS: {
            const template = fromJS(action.template);
            const accountId = template.get('accountId');
            const index = state.getIn([accountId, 'templates', 'items']).findIndex(t => t.get('id') === template.get('id'));
            if (index !== -1) {
                return state.setIn([accountId, 'templates', 'items', index], template);
            }
            return state.updateIn([accountId, 'templates', 'items'], items => items.push(template));
        }

        case FETCH_RECURRING_TRANSACTIONS_START:
            return state.setIn([action.accountId, 'recurringTransactions', 'fetching'], true);

        case RECEIVE_RECURRING_TRANSACTIONS_SUCCESS:
            return state.setIn([action.accountId, 'recurringTransactions'], Map({
                items: fromJS(action.recurringTransactions).map(r => r.updateIn(['transaction', 'date'], d => dateStringToDate(d))),
                fetched: true,
                fetching: false
            }));

        case EDIT_RECURRING_TRANSACTION_SUCCESS: {
            let recurringTransaction = fromJS(action.recurringTransaction);
            recurringTransaction = recurringTransaction.updateIn(['transaction', 'date'], d => dateStringToDate(d));
            const accountId = recurringTransaction.getIn(['transaction', 'accountId']);
            const index = state.getIn([accountId, 'recurringTransactions', 'items']).findIndex(t => t.get('id') === recurringTransaction.get('id'));
            if (index !== -1) {
                return state.setIn([accountId, 'recurringTransactions', 'items', index], recurringTransaction);
            }
            return state.updateIn([accountId, 'recurringTransactions', 'items'], items => items.push(recurringTransaction));
        }

        case DELETE_RECURRING_TRANSACTION_SUCCESS: {
            const recurringTransaction = action.recurringTransaction;
            const accountId = recurringTransaction.transaction.accountId;
            return state.updateIn([accountId, 'recurringTransactions', 'items'], items => items.filter(rt => rt.get('id') !== recurringTransaction.id));
        }

        default:
            return state;
    }
}

function calculateValues(data, futureValue) {
    let ret = List();
    data
        .valueSeq()
        .sortBy(t => t.get('date'))
        .reverse()
        .forEach(t => {
            ret = ret.push(t.set('balance', futureValue));
            futureValue -= t.get('amount');
        });
    return ret;
}
