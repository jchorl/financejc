import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { Map } from 'immutable';
import { connect } from 'react-redux';
import TransactionWrapper from './TransactionWrapper';
import TransactionEdit from '../TransactionEdit';
import { fetchTemplates, fetchTransactions } from '../../actions/accountData';
import { NEW_TRANSACTION_ID } from '../../constants';
import './TransactionList.css';

export function emptyTransaction(accountId) {
    return Map({
        id: NEW_TRANSACTION_ID,
        name: '',
        date: new Date(),
        category: '',
        amount: 0,
        note: '',
        accountId
    });
}

class TransactionList extends Component {
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
                                 transactions: ImmutablePropTypes.contains({
                                     fetched: PropTypes.bool.isRequired,
                                     fetching: PropTypes.bool.isRequired,
                                     calculatedValues: ImmutablePropTypes.list
                                 }).isRequired,
                                 templates: ImmutablePropTypes.contains({
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
                    }).isRequired
                    ).isRequired
        }).isRequired
    }

    constructor() {
        super();

        this.state = {
            isEnteringTransaction: false
        }
    }

    fetchTransactionsIfNotFetched = (accountId, props) => {
        const {
            accountData,
            dispatch,
        } = props;

        if (!accountData.getIn([accountId, 'transactions', 'fetched']) && !accountData.getIn([accountId, 'transactions', 'fetching'])) {
            dispatch(fetchTransactions(accountId));
        }
    }

    fetchTemplatesIfNotFetched = (accountId, props) => {
        const {
            accountData,
            dispatch,
        } = props;

        if (!accountData.getIn([accountId, 'templates', 'fetched']) && !accountData.getIn([accountId, 'templates', 'fetching'])) {
            dispatch(fetchTemplates(accountId));
        }
    }

    componentWillMount() {
        const {
            match: { params: { id } }
        } = this.props;

        let accountId = parseInt(id, 10);
        this.fetchTransactionsIfNotFetched(accountId, this.props);
        this.fetchTemplatesIfNotFetched(accountId, this.props);
    }

    componentWillReceiveProps(nextProps) {
        const {
            match: { params: { id } }
        } = nextProps;

        let accountId = parseInt(id, 10);
        this.fetchTransactionsIfNotFetched(accountId, nextProps);
        this.fetchTemplatesIfNotFetched(accountId, nextProps);
    }

    newTransaction = transactionTemplate => () => {
        this.setState({
            isEnteringTransaction: true,
            transactionTemplate
        });
    }

    doneEnteringTransaction = () => {
        this.setState({
            isEnteringTransaction: false
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
            isEnteringTransaction,
            transactionTemplate
        } = this.state;

        const parsedId = parseInt(id, 10);
        const accountCurrencyCode = account.getIn(['items', parsedId, 'currency']);
        const accountCurrency = currency.getIn(['items', accountCurrencyCode]);

        return (
                <div className="transactionList">
                    {
                    isEnteringTransaction
                    ? (
                    <TransactionEdit
                        transaction={ transactionTemplate }
                        currency={ accountCurrency }
                        exitEditMode={ this.doneEnteringTransaction }
                    />
                    )
                    : (
                    <div className="templateBar">
                        <button key="EMPTY" className="templateButton" onClick={ this.newTransaction(emptyTransaction(parsedId)) }><i className="fa fa-plus-circle"></i> Empty</button>
                        {
                        accountData.getIn([parsedId, 'templates', 'items']).map(
                        t => <button key={ t.get('id') } className="templateButton" onClick={ this.newTransaction(t) }>{ t.get('templateName') }</button>
                        )
                        }
                    </div>
                    )
                    }
                    {
                    accountData.getIn([parsedId, 'transactions', 'calculatedValues']).map(
                    t => <TransactionWrapper key={ t.get('id') } transaction={ t } currency={ accountCurrency } />
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
})(TransactionList);
