import React from 'react';
import { connect } from 'react-redux';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { fetchTemplates } from '../../actions';
import styles from './templateList.css';
import { Template, TemplateForm } from '../template';


function fetchTemplatesIfNecessary(props) {
    const {
        accountId,
        accountTemplate,
        dispatch
    } = props;

    if (!accountTemplate.get(accountId).get('fetched')) {
        dispatch(fetchTemplates(accountId));
    }
}

@connect((state) => {
    return {
        accountTemplate: state.accountTemplate
    };
})
export default class TemplateList extends React.Component {
    static propTypes = {
        accountId: React.PropTypes.number.isRequired,
        accountTemplate: ImmutablePropTypes.map.isRequired,
        currency: ImmutablePropTypes.map.isRequired,
        dispatch: React.PropTypes.func.isRequired
    };

    constructor (props) {
        super(props);
        this.state = {
            newTemplate: false
        };

        fetchTemplatesIfNecessary(props);
    }

    componentWillReceiveProps(newProps) {
        fetchTemplatesIfNecessary(newProps);
    }

    startNewTemplate = () => {
        this.setState({
            newTemplate: true
        });
    }

    exitNewTemplate = () => {
        this.setState({
            newTemplate: false
        });
    }

    render () {
        const {
            accountId,
            accountTemplate,
            currency
        } = this.props;

        let transactions = accountTemplate.get(accountId).get('templates');

        return (
            <div>
                <div className={ styles.headings }>
                    <span className={ styles.column }>Template Name</span>
                    <span className={ styles.column }>Name</span>
                    <span className={ styles.column }>Category</span>
                    <span className={ styles.column }>Amount</span>
                </div>
                { !this.state.newTemplate ?
                        (
                            <button className={ styles.newTemplate } onClick={ this.startNewTemplate }>
                                New
                            </button>
                        ) : (
                            <TemplateForm accountId={ accountId } form='new' done={ this.exitNewTemplate } currency={ currency } initialValues={ { templateName: '', name: '', category: '', amount: 0 } } />
                        )
                }
                { transactions.map(transaction => (<Template key={ transaction.get('id') } template={ transaction } currency={ currency }/>)).toOrderedSet().toArray() }
            </div>
        );
    }
}
