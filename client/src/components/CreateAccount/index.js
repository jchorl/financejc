import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { createAccount } from '../../actions/account';
import './CreateAccount.css';

class CreateAccount extends Component {
    static propTypes = {
        currency: ImmutablePropTypes.contains({
            items: ImmutablePropTypes.mapOf(
                           ImmutablePropTypes.contains({
                               code: PropTypes.string.isRequired,
                               name: PropTypes.string.isRequired
                           })
                           )
        })
    }

    constructor() {
        super();

        this.state = {
            currencyVal: 'USD'
        }
    }

    handleChange = e => {
        this.setState({ currencyVal: e.target.value });
    }

    createAccount = e => {
        const {
            dispatch,
            history
        } = this.props;

        const account = {
            name: e.target.name.value,
            currency: e.target.currency.value
        }

        dispatch(createAccount(account))
            .then(a => history.push(`/accounts/${a.id}/transactions`));
        e.preventDefault();
    }

    render() {
        const { currency } = this.props;

        return (
                <div className="container createAccount">
                    <h1>Create New Account</h1>
                    <div className="divider"></div>
                    <form className="createAccountForm" onSubmit={ this.createAccount }>
                        <label htmlFor="name">Name:</label><input type="text" placeholder="Schwab Checking" name="name" />
                        <label htmlFor="currency">Currency:</label>
                        <select name="currency" value={ this.state.currencyVal } onChange={ this.handleChange }>
                            {
                            currency.get('items')
                            .valueSeq()
                            .sortBy(c => c.get('code'))
                            .map(c => <option key={ c.get('code') } value={ c.get('code') }>{ `${c.get('code')} (${c.get('name')})` }</option>)
                            }
                        </select>
                        <button className="submitButton" type="submit">Create</button>
                    </form>
                </div>
                );
    }
}

export default withRouter(connect(state => {
    return {
        currency: state.currency
    };
})(CreateAccount));
