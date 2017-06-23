import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { CREATE_ACCOUNT_ID, SUMMARY_ID } from '../../constants';
import AccountEntry from './AccountEntry';
import CreateAccountEntry from './CreateAccountEntry';
import SummaryEntry from './SummaryEntry';
import './AccountsBar.css';

class AccountsBar extends Component {
    static propTypes = {
        account: ImmutablePropTypes.contains({
            items: ImmutablePropTypes.mapOf(
                           ImmutablePropTypes.contains({
                               id: PropTypes.number.isRequired
                           }).isRequired
                           ).isRequired
        }).isRequired,
        selected: PropTypes.string.isRequired
    }

    render() {
        const {
            account,
            selected
        } = this.props;

        return (
                <div className="container accountsBar">
                    <div className="accountsLabel">
                        <h1>Accounts</h1>
                        <div className="divider"></div>
                    </div>
                    <div className="entries">
                        {
                        !account.get('items').isEmpty()
                        ? (
                        <div>
                            <SummaryEntry selected={ selected === SUMMARY_ID } />
                            <div className="divider"></div>
                        </div>
                        )
                        : null
                        }
                        {
                        account
                        .get('items')
                        .valueSeq()
                        .sortBy(a => a.get('name'))
                        .map(a =>
                        <div key={ a.get('id') }>
                            <AccountEntry account={ a } selected={ selected === a.get('id') + '' } />
                            <div className="divider"></div>
                        </div>
                        )
                        }
                        <CreateAccountEntry selected={ selected === CREATE_ACCOUNT_ID } />
                    </div>
                </div>
                );
    }
}

export default connect(state => {
    return {
        account: state.account
    }
})(AccountsBar);
