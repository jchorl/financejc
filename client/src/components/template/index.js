import React from 'react';
import classNames from 'classnames';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import Autosuggest from 'react-autosuggest';
import styles from './template.css';
import { toCurrency, toDecimal, toWhole, queryByFieldAndVal } from '../../utils';
import { putTemplate, deleteTemplate } from '../../actions';

export class Template extends React.Component {
    static propTypes = {
        template: React.PropTypes.object,
        currency: ImmutablePropTypes.map.isRequired,
    };

    constructor(props) {
        super(props);
        this.state = { editMode: false };
    }

    enterEditMode = () => {
        this.setState({ editMode: true });
    }

    exitEditMode = () => {
        this.setState({ editMode: false });
    }

    save = () => {
        this.exitEditMode();
    }

    render () {
        const {
            template,
            currency
        } = this.props;

        return template ? (
            this.state.editMode ? (
                <TemplateForm template={ template } initialValues={ getFormInitialValues(template, currency) } done={ this.exitEditMode } currency={ currency } />
            ) : (
                <div className={ classNames(styles.template, styles.templateFields) }>
                    <span className={ classNames(styles.templateField, styles.nonEdit) } onClick={ this.enterEditMode }>{ template.get('templateName') }</span>
                    <span className={ classNames(styles.templateField, styles.nonEdit) } onClick={ this.enterEditMode }>{ template.get('name') }</span>
                    <span className={ classNames(styles.templateField, styles.nonEdit) } onClick={ this.enterEditMode }>{ template.get('category') }</span>
                    <span className={ classNames(styles.templateField, styles.nonEdit) } onClick={ this.enterEditMode }>{ toCurrency(toDecimal(template.get('amount'), currency.get('digitsAfterDecimal')), currency.get('code')) }</span>
                </div>
            )
        ) : null;
    }
}

function getFormInitialValues(template, currency) {
    return {
        templateName: template.get('templateName'),
        name: template.get('name'),
        category: template.get('category'),
        amount: toDecimal(template.get('amount'), currency.get('digitsAfterDecimal'))
    };
}

function getSuggestionValue(field, suggestion) {
    return suggestion[field];
}

function renderSuggestion(field, suggestion) {
    return (
        <span>{ suggestion[field] }</span>
    );
}

@connect()
export class TemplateForm extends React.Component {
    static propTypes = {
        currency: ImmutablePropTypes.map.isRequired,
        dispatch: React.PropTypes.func.isRequired,
        initialValues: React.PropTypes.object.isRequired,
        // either template (for editing) or accountId (for new templates) should be passed
        accountId: React.PropTypes.number,
        template: ImmutablePropTypes.map,
        done: React.PropTypes.func
    };

    constructor(props) {
        super(props);

        this.state = {
            values: props.initialValues,
            suggestions: {
                name: [],
                category: []
            }
        };

        this.lastRequestId = null;
    }

    componentWillReceiveProps(nextProps) {
        this.setState({
            values: nextProps.initialValues
        });
    }

    loadSuggestions = (field, value) => {
        const {
            accountId,
            template
        } = this.props;

        let resolvedAccountId = accountId;
        if (template) {
            resolvedAccountId = template.get('accountId');
        }

        let id = Math.random();
        this.setState({
            lastRequestId: id
        });

        let that = this;

        // ideally requests are made from actions, buuuuut it is much easier and faster to skip redux
        queryByFieldAndVal(resolvedAccountId, field, value).then(templates => {
            // stale query
            if (id !== that.state.lastRequestId) {
                return;
            }

            let newState = { suggestions: {} };
            Object.assign(newState.suggestions, that.state.suggestions);
            newState.suggestions[field] = templates;
            that.setState(newState);
        });
    }

    onChange = field => {
        let that = this;
        return function(event, { newValue }) {
            let values = Object.assign({}, that.state.values);
            values[field] = newValue;
            that.setState({
                values
            });
        };
    };

    onSuggestionsFetchRequested = (field, { value }) => {
        this.loadSuggestions(field, value);
    };

    onSuggestionsClearRequested = () => {
        this.setState({
            suggestions: {
                name: [],
                category: []
            }
        });
    };

    onNameSuggestionSelected = (event, { suggestion }) => {
        const { currency } = this.props;

        this.setState({
            values: {
                name: suggestion.name,
                category: suggestion.category,
                amount: toDecimal(suggestion.amount, currency.get('digitsAfterDecimal'))
            }
        });
    }

    fieldChange = name => e => {
        let newState = Object.assign({}, this.state.values);
        newState[name] = e.target.value;
        this.setState({
            values: newState
        });
    }

    deleteTemplate = (e) => {
        const {
            dispatch,
            done,
            template
        } = this.props;

        dispatch(deleteTemplate(template.get('id'), template.get('accountId')));
        done && done();
        e.preventDefault();
    }

    submit = (e) => {
        const {
            accountId,
            currency,
            dispatch,
            done,
            template
        } = this.props;

        let obj = {
            templateName: e.target['templateName'].value,
            name: e.target['name'].value,
            category: e.target['category'].value,
            amount: toWhole(parseFloat(e.target['amount'].value), currency.get('digitsAfterDecimal')),
            accountId
        };

        if (template) {
            obj = Object.assign(template.toObject(), obj);
            obj.accountId = template.get('accountId');
        }

        dispatch(putTemplate(obj));
        done && done();
        e.preventDefault();
    }

    render () {
        const {
            done,
            template
        } = this.props;

        const {
            suggestions,
            values
        } = this.state;

        const nameInputProps = {
            name: 'name',
            placeholder: 'Name',
            value: values.name,
            onChange: this.onChange('name')
        };

        const categoryInputProps = {
            name: 'category',
            placeholder: 'Category',
            value: values.category,
            onChange: this.onChange('category')
        };

        return (
            <div className={ styles.template }>
                <form onSubmit={ this.submit }>
                    <div className={ styles.templateFields }>
                        <input type="text" name="templateName" placeholder="Template Name" value={ values.templateName } onChange={ this.fieldChange('templateName') } className={ styles.templateField } />
                        <Autosuggest
                            id="name"
                            suggestions={ suggestions.name }
                            onSuggestionsFetchRequested={ this.onSuggestionsFetchRequested.bind(this, 'name') }
                            onSuggestionsClearRequested={ this.onSuggestionsClearRequested }
                            onSuggestionSelected={ this.onNameSuggestionSelected }
                            getSuggestionValue={ getSuggestionValue.bind(undefined, 'name') }
                            renderSuggestion={ renderSuggestion.bind(undefined, 'name') }
                            inputProps={ nameInputProps }
                            theme={ styles } />
                        <Autosuggest
                            id="category"
                            suggestions={ suggestions.category }
                            onSuggestionsFetchRequested={ this.onSuggestionsFetchRequested.bind(this, 'category') }
                            onSuggestionsClearRequested={ this.onSuggestionsClearRequested }
                            getSuggestionValue={ getSuggestionValue.bind(undefined, 'category') }
                            renderSuggestion={ renderSuggestion.bind(undefined, 'category') }
                            inputProps={ categoryInputProps }
                            theme={ styles } />
                        <input type="text" name="amount" value={ values.amount } onChange={ this.fieldChange('amount') } placeholder="0" className={ styles.templateField } />
                    </div>
                    <div className={ styles.saveExit }>
                        <button type="button" onClick={ done }>Cancel</button>
                        {
                            template ? <button type="button" onClick={ this.deleteTemplate }>Delete</button> : null
                        }
                        <button type="submit">Save</button>
                    </div>
                </form>
            </div>
        );
    }
}
