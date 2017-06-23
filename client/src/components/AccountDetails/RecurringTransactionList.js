import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { Map } from 'immutable';
import { NEW_RECURRING_TRANSACTION_ID, SCHEDULE_TYPES } from '../../constants';
import { emptyTransaction } from './TransactionList';
import RecurringTransactionWrapper from './RecurringTransactionWrapper';
import { fetchRecurringTransactions } from '../../actions/accountData';
import RecurringTransactionEdit from './RecurringTransactionEdit.js';
import './RecurringTransactionList.css';

function emptyRecurringTransaction(accountId) {
    return Map({
        id: NEW_RECURRING_TRANSACTION_ID,
        transaction: emptyTransaction(accountId),
        scheduleType: SCHEDULE_TYPES.FIXED_DAY_MONTH,
        secondsBetween: null,
        dayOf: 1,
        secondsBeforeToPost: 172800
    });
}

class RecurringTransactionList extends Component {
    static propTypes = {
        account: ImmutablePropTypes.contains({
            items: ImmutablePropTypes.mapOf(
                           ImmutablePropTypes.contains({
                               id: PropTypes.number.isRequired,
                               currency: PropTypes.string.isRequired
                           }).isRequired
                           ).isRequired
        }).isRequired,
        accountData: ImmutablePropTypes.mapOf(
                             ImmutablePropTypes.contains({
                                 recurringTransactions: ImmutablePropTypes.contains({
                                     fetched: PropTypes.bool.isRequired,
                                     fetching: PropTypes.bool.isRequired,
                                     items: ImmutablePropTypes.list
                                 }).isRequired
                             }).isRequired
                             ).isRequired,
        currency: ImmutablePropTypes.contains({
            items: ImmutablePropTypes.mapOf(
                    ImmutablePropTypes.contains({
                        digitsAfterDecimal: PropTypes.number.isRequired
                    })
                    ).isRequired
        }).isRequired
    }

    constructor(props) {
        super(props);
        this.state = {
            isEnteringRecurringTransaction: false
        };
    }

    fetchRecurringTransactionsIfNotFetched = (accountId, props) => {
        const {
            accountData,
            dispatch,
        } = props;

        if (!accountData.getIn([accountId, 'recurringTransactions', 'fetched']) && !accountData.getIn([accountId, 'recurringTransactions', 'fetching'])) {
            dispatch(fetchRecurringTransactions(accountId));
        }
    }

    componentWillMount() {
        const {
            match: { params: { id } }
        } = this.props;

        let accountId = parseInt(id, 10);
        this.fetchRecurringTransactionsIfNotFetched(accountId, this.props);
    }

    componentWillReceiveProps(nextProps) {
        const {
            match: { params: { id } }
        } = nextProps;

        let accountId = parseInt(id, 10);
        this.fetchRecurringTransactionsIfNotFetched(accountId, nextProps);
    }

    newRecurringTransaction = recurringTransaction => () => {
        recurringTransaction = recurringTransaction.set('id', NEW_RECURRING_TRANSACTION_ID);
        this.setState({
            isEnteringRecurringTransaction: true,
            recurringTransaction
        });
    }

    doneEnteringRecurringTransaction = () => {
        this.setState({
            isEnteringRecurringTransaction: false
        });
    }

    render() {
        const {
            account,
            accountData,
            currency,
            match: { params: { id } }
        } = this.props;

        const {
            isEnteringRecurringTransaction,
            recurringTransaction
        } = this.state;

        const parsedId = parseInt(id, 10);
        const accountCurrencyCode = account.getIn(['items', parsedId, 'currency']);
        const accountCurrency = currency.getIn(['items', accountCurrencyCode]);

        return (
                <div className="recurringTransactionList">
                    {
                    isEnteringRecurringTransaction
                    ? (
                    <RecurringTransactionEdit
                        recurringTransaction={ recurringTransaction }
                        currency={ accountCurrency }
                        exitEditMode={ this.doneEnteringRecurringTransaction }
                    />
                    )
                    : (
                    <div className="newRecurringTransactionBar">
                        <button key="EMPTY" onClick={ this.newRecurringTransaction(emptyRecurringTransaction(parsedId)) }><i className="fa fa-plus-circle"></i> New Recurring Transaction</button>
                    </div>
                    )
                    }
                    {
                    accountData.getIn([parsedId, 'recurringTransactions', 'items']).map(
                    t => <RecurringTransactionWrapper key={ t.get('id') } recurringTransaction={ t } currency={ accountCurrency } />
                    )
                    }
                </div>
                );
    }
}

export default connect(state => {
    return {
        account: state.account,
        accountData: state.accountData,
        currency: state.currency
    }
})(RecurringTransactionList);
