import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { toCurrencyString } from '../../util';

export default class Template extends Component {
    static propTypes = {
        template: ImmutablePropTypes.contains({
            templateName: PropTypes.string.isRequired,
            name: PropTypes.string.isRequired,
            category: PropTypes.string.isRequired,
            amount: PropTypes.number.isRequired
        }).isRequired,
        currency: ImmutablePropTypes.contains({
            code: PropTypes.string.isRequired,
            digitsAfterDecimal: PropTypes.number.isRequired
        }).isRequired,
        enterEditMode: PropTypes.func.isRequired
    }

    render() {
        const {
            template,
            currency,
            enterEditMode
        } = this.props;

        return (
                <div className="template" onClick={ enterEditMode } >
                    <div>{ template.get('templateName') }</div>
                    <div>{ template.get('name') }</div>
                    <div>{ template.get('category') }</div>
                    <div>{ toCurrencyString(template.get('amount'), currency.get('code'), currency.get('digitsAfterDecimal')) }</div>
                </div>
                );
    }
}
