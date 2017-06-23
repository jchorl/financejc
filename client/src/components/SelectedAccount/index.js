import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Route, Switch } from 'react-router-dom';
import AccountsBar from '../AccountsBar';
import Summary from '../Summary';
import AccountDetails from '../AccountDetails';
import CreateAccount from '../CreateAccount';
import './SelectedAccount.css';

export default class SelectedAccount extends Component {
    static propTypes = {
        match: PropTypes.shape({
            params: PropTypes.shape({
                id: PropTypes.string.isRequired
            }).isRequired
        }).isRequired
    }

    render() {
        const { match: { params: { id } } } = this.props;

        return (
                <div className="selectedAccount">
                    <div className="accountsBarWrapper">
                        <AccountsBar selected={ id } />
                    </div>
                    <div className="restWrapper">
                        <Switch>
                            <Route path="/accounts/create" component={ CreateAccount } />
                            <Route path="/accounts/summary" component={ Summary } />
                            <Route path="/accounts/:id/:page" component={ AccountDetails } />
                        </Switch>
                    </div>
                </div>
                );
    }
}
