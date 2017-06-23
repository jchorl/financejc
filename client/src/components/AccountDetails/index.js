import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { Route, withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import classNames from 'classnames';
import { ACCOUNT_PAGES } from '../../constants';
import TransactionList from './TransactionList';
import TemplateList from './TemplateList';
import RecurringTransactionList from './RecurringTransactionList';
import './AccountDetails.css';

class AccountDetails extends Component {
    static propTypes = {
        account: ImmutablePropTypes.contains({
            items: ImmutablePropTypes.mapOf(
                           ImmutablePropTypes.contains({
                               id: PropTypes.number.isRequired,
                               name: PropTypes.string.isRequired
                           }).isRequired
                           ).isRequired
        }).isRequired,
        match: PropTypes.shape({
            params: PropTypes.shape({
                id: PropTypes.string.isRequired,
                page: PropTypes.string.isRequired
            }).isRequired
        }).isRequired,
        history: PropTypes.shape({
            push: PropTypes.func.isRequired
        })
    }

    goTo = page => () => {
        const {
            history,
            match: { params: { id } }
        } = this.props;

        history.push(`/accounts/${id}/${page}`);
    }

    render() {
        const {
            account,
            match: { params: { id, page } }
        } = this.props;

        const parsedId = parseInt(id, 10);
        const acc = account.getIn(['items', parsedId]);

        return acc
            ? (
                    <div className="accountDetails container">
                        <div className="nameAndNav">
                            <h1>{ acc.get('name') }</h1>
                            <div className="accountNav">
                                <div className={ classNames('accountNavItem', { selected: page === ACCOUNT_PAGES.TRANSACTIONS }) } onClick={ this.goTo('transactions') }>
                                    <div className="label">Transactions</div>
                                    <div className="underline"></div>
                                </div>
                                <div className={ classNames('accountNavItem', { selected: page === ACCOUNT_PAGES.TEMPLATES }) } onClick={ this.goTo('templates') }>
                                    <div className="label">Templates</div>
                                    <div className="underline"></div>
                                </div>
                                <div className={ classNames('accountNavItem', { selected: page === ACCOUNT_PAGES.RECURRING_TRANSACTIONS }) } onClick={ this.goTo('recurringTransactions') }>
                                    <div className="label">Recurring Transactions</div>
                                    <div className="underline"></div>
                                </div>
                            </div>
                        </div>
                        <div className="divider"></div>
                        <Route path="/accounts/:id/transactions" component={ TransactionList } />
                        <Route path="/accounts/:id/templates" component={ TemplateList } />
                        <Route path="/accounts/:id/recurringTransactions" component={ RecurringTransactionList } />
                    </div>
                    )
            : null;
    }
}

export default withRouter(connect(state => {
    return {
        account: state.account
    }
})(AccountDetails));
