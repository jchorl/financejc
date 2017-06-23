import React, { Component } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import RecurringTransaction from './RecurringTransaction';
import RecurringTransactionEdit from './RecurringTransactionEdit';

export default class RecurringTransactionWrapper extends Component {
    static propTypes = {
        recurringTransaction: ImmutablePropTypes.map.isRequired,
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
            recurringTransaction,
            currency
        } = this.props;

        const { isEditMode } = this.state;

        return isEditMode
            ? (
                    <RecurringTransactionEdit recurringTransaction={ recurringTransaction } currency={ currency } exitEditMode={ this.exitEditMode } />
                    )
            : (
                    <RecurringTransaction recurringTransaction={ recurringTransaction } currency={ currency } enterEditMode={ this.enterEditMode } />
                    );
    }
}
