import React from 'react';
import { connect } from 'react-redux';
import ImmutablePropTypes from 'react-immutable-proptypes';
import classNames from 'classnames';
import { toCurrency, toDecimal } from '../../utils';
import { selectAccount } from '../../actions';
import styles from './accountList.css';

@connect((state) => {
    return {
        accounts: state.accounts,
        currencies: state.currencies
    };
})
export default class AccountList extends React.Component {
    static propTypes = {
        accounts: ImmutablePropTypes.map.isRequired,
        currencies: ImmutablePropTypes.map.isRequired,
        dispatch: React.PropTypes.func.isRequired
    }

    render () {
        const {
            accounts,
            currencies,
            dispatch
        } = this.props;

        let selectedClass = {};
        selectedClass[styles.selected] = accounts.get('selected') === -1;
        return (
            <div>
                <h3 className={ styles.accountsTitle }>Accounts</h3>
                <div>
                    <button key={ 'new' } className={ classNames(styles.accountButton, selectedClass) } onClick={ dispatch.bind(undefined, selectAccount(-1)) }>
                        <div className={ styles.accountName }>
                            New
                        </div>
                    </button>
                    { accounts.get('accounts').map(account => {
                        let selectedClass = {};
                        selectedClass[styles.selected] = accounts.get('selected') === account.get('id');
                        return (
                            <button key={ account.get('id') } className={ classNames(styles.accountButton, selectedClass) } onClick={ dispatch.bind(undefined, selectAccount(account.get('id'))) }>
                                <div className={ styles.accountName }>
                                    { account.get('name') }
                                </div>
                                <div className={ styles.accountInfo }>
                                    Balance: { toCurrency(toDecimal(account.get('futureValue'), currencies.get('currencies').get(account.get('currency')).get('digitsAfterDecimal')), account.get('currency')) }
                                </div>
                            </button>
                        );
                    }).valueSeq().toArray() }
                </div>
            </div>
        );
    }
}
