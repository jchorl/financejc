import React, { Component } from 'react';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import classNames from 'classnames';
import { ACCOUNT_PAGES } from '../../constants';
import { toCurrencyString } from '../../util';

class AccountEntry extends Component {
    static propTypes = {
        account: ImmutablePropTypes.contains({
                               id: PropTypes.number.isRequired,
                               futureValue: PropTypes.number.isRequired,
                               currency: PropTypes.string.isRequired,
                               name: PropTypes.string.isRequired
                           }).isRequired,
        currency: ImmutablePropTypes.contains({
            items: ImmutablePropTypes.mapOf(
                    ImmutablePropTypes.contains({
                        digitsAfterDecimal: PropTypes.number.isRequired
                    })
                    )
        }).isRequired,
        history: PropTypes.shape({
            push: PropTypes.func.isRequired
        }),
        selected: PropTypes.bool.isRequired
    }

    goToAccount = () => {
        const {
            account,
            history
        } = this.props;

        history.push(`/accounts/${account.get('id')}/${ACCOUNT_PAGES.TRANSACTIONS}`);
    }

    render() {
        const {
            account,
            currency,
            selected
        } = this.props;

        const accountCurrency = account.get('currency');

        return (
                <div className="entry accountEntry" onClick={ this.goToAccount }>
                    <div className={ classNames("selectedBlock", { selected }) }></div>
                    <div className="name">{ account.get('name') }</div>
                    <div className="currency">{ accountCurrency }</div>
                    <div className="value">{ toCurrencyString(account.get('futureValue'), accountCurrency, currency.get('items').get(accountCurrency).get('digitsAfterDecimal')) }</div>
                </div>
                );
    }
}

export default withRouter(connect(state => {
    return {
        currency: state.currency
    }
})(AccountEntry));
