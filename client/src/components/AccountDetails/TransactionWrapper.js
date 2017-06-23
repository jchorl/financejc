import React, { Component } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import Transaction from './Transaction';
import TransactionEdit from '../TransactionEdit';

export default class TransactionWrapper extends Component {
    static propTypes = {
        transaction: ImmutablePropTypes.map.isRequired,
        currency: ImmutablePropTypes.map.isRequired
    }

    constructor() {
        super();
        this.state = { isEditMode: false }
    }

    enterEditMode = () => {
        this.setState({ isEditMode: true });
    }

    exitEditMode = () => {
        this.setState({ isEditMode: false });
    }

    render() {
        const {
            transaction,
            currency
        } = this.props;

        const { isEditMode } = this.state;

        return isEditMode
            ? (
                    <TransactionEdit transaction={ transaction } currency={ currency } exitEditMode={ this.exitEditMode } />
                    )
            : (
                    <Transaction transaction={ transaction } currency={ currency } enterEditMode={ this.enterEditMode } />
                    );
    }
}
