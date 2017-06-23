import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { NEW_TEMPLATE_ID } from '../../constants';
import { editTemplate, newTemplate } from '../../actions/accountData';
import SuggestedInput from '../SuggestedInput';

class TemplateEdit extends Component {
    static propTypes = {
        template: ImmutablePropTypes.contains({
            id: PropTypes.number.isRequired,
            templateName: PropTypes.string.isRequired,
            name: PropTypes.string.isRequired,
            category: PropTypes.string.isRequired,
            amount: PropTypes.number.isRequired,
            accountId: PropTypes.number.isRequired
        }).isRequired,
        currency: ImmutablePropTypes.contains({
            code: PropTypes.string.isRequired,
            digitsAfterDecimal: PropTypes.number.isRequired
        }).isRequired,
        dispatch: PropTypes.func.isRequired,
        exitEditMode: PropTypes.func.isRequired
    }

    constructor(props) {
        super(props);

        const {
            template
        } = this.props;

        this.state = {
            name: template.get('name'),
            templateName: template.get('templateName'),
            category: template.get('category'),
            amount: template.get('amount'),
            savedName: '',
            savedCategory: '',
            savedAmount: 0,
            isSaved: false
        }
    }

    handleChange = field => e => {
        let newState = {};
        switch (field) {
            case 'amount': {
                const { currency } = this.props;

                newState[field] = parseFloat(e.target.value, 10) * Math.pow(10, currency.get('digitsAfterDecimal'));
                break;
            }

            default:
                newState[field] = e.target.value;
        }

        this.setState(newState);
    }

    onSubmit = e => {
        const {
            dispatch,
            exitEditMode,
            template: originalTemplate
        } = this.props;

        const {
            templateName,
            name,
            category,
            amount
        } = this.state;

        const template = {
            templateName,
            name,
            category,
            amount: Math.round(amount),
            accountId: originalTemplate.get('accountId')
        };

        if (originalTemplate.get('id') === NEW_TEMPLATE_ID) {
            dispatch(newTemplate(template))
                .then(exitEditMode);
        } else {
            template.id = originalTemplate.get('id');
            dispatch(editTemplate(template))
                .then(exitEditMode);
        }

        e.preventDefault();
    }

    selectNameSuggestion = suggestion => {
        this.setState({
            name: suggestion.get('name'),
            category: suggestion.get('category'),
            amount: suggestion.get('amount')
        });
    }

    hoverNameSuggestion = suggestion => {
        const {
            name,
            category,
            amount,
            isSaved
        } = this.state;

        let newState = {
            name: suggestion.get('name'),
            category: suggestion.get('category'),
            amount: suggestion.get('amount')
        };

        if (!isSaved) {
            newState.savedName = name;
            newState.savedCategory = category;
            newState.savedAmount = amount;
            newState.isSaved = true;
        }

        this.setState(newState);
    }

    unhoverNameSuggestion = () => {
        const {
            savedName,
            savedCategory,
            savedAmount
        } = this.state;

        this.setState({
            name: savedName,
            category: savedCategory,
            amount: savedAmount,
            isSaved: false
        });
    }

    selectCategorySuggestion = suggestion => {
        this.setState({
            category: suggestion.get('category')
        });
    }

    hoverCategorySuggestion = suggestion => {
        const {
            category,
            isSaved
        } = this.state;

        let newState = {
            category: suggestion.get('category')
        };

        if (!isSaved) {
            newState.savedCategory = category;
            newState.isSaved = true;
        }

        this.setState(newState);
    }

    unhoverCategorySuggestion = () => {
        const {
            savedCategory
        } = this.state;

        this.setState({
            category: savedCategory,
            isSaved: false
        });
    }

    render() {
        const {
            template,
            currency,
            exitEditMode
        } = this.props;

        return (
                <form className="templateEdit" onSubmit={ this.onSubmit }>
                    <div className="inputWithLabel">
                        <label htmlFor="templateName">Template Name: </label>
                        <input
                            type="text"
                            name="templateName"
                            value={ this.state.templateName }
                            onChange={ this.handleChange('templateName') }
                        />
                    </div>
                    <div className="inputWithLabel">
                        <label htmlFor="name">Name: </label>
                        <SuggestedInput
                            type="text"
                            name="name"
                            autoComplete="off"
                            value={ this.state.name }
                            onChange={ this.handleChange('name') }
                            accountId={ template.get('accountId') }
                            field="name"
                            selectSuggestion={ this.selectNameSuggestion }
                            hoverSuggestion={ this.hoverNameSuggestion }
                            unhoverSuggestion={ this.unhoverNameSuggestion }
                        />
                    </div>
                    <div className="inputWithLabel">
                        <label htmlFor="category">Category: </label>
                        <SuggestedInput
                            type="text"
                            name="category"
                            autoComplete="off"
                            value={ this.state.category }
                            onChange={ this.handleChange('category') }
                            accountId={ template.get('accountId') }
                            field="category"
                            selectSuggestion={ this.selectCategorySuggestion }
                            hoverSuggestion={ this.hoverCategorySuggestion }
                            unhoverSuggestion={ this.unhoverCategorySuggestion }
                        />
                    </div>
                    <div className="inputWithLabel">
                        <label htmlFor="amount">Amount: </label>
                        <input
                            type="number"
                            step="0.01"
                            name="amount"
                            value={ this.state.amount / Math.pow(10, currency.get('digitsAfterDecimal')) }
                            onChange={ this.handleChange('amount') }
                        />
                    </div>
                    <div>
                        <button type="button" onClick={ exitEditMode }>Cancel</button>
                        <button type="submit">Save</button>
                    </div>
                </form>
                );
    }
}

export default connect()(TemplateEdit);
