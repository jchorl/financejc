import React from 'react';
import { connect } from 'react-redux';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { fetchTransactionTemplates } from '../../actions';
import styles from './transactionTemplateList.css';
import { TransactionTemplate, TransactionTemplateForm } from '../transactionTemplate';


function fetchTransactionTemplatesIfNecessary(props) {
    const {
        accountId,
        accountTransactionTemplate,
        dispatch
    } = props;

    if (!accountTransactionTemplate.get(accountId).get('fetched')) {
        dispatch(fetchTransactionTemplates(accountId));
    }
}

@connect((state) => {
    return {
        accountTransactionTemplate: state.accountTransactionTemplate
    };
})
export default class TransactionTemplateList extends React.Component {
    static propTypes = {
        accountId: React.PropTypes.number.isRequired,
        accountTransactionTemplate: ImmutablePropTypes.map.isRequired,
        currency: ImmutablePropTypes.map.isRequired,
        dispatch: React.PropTypes.func.isRequired
    };

    constructor (props) {
        super(props);
        this.state = {
            newTransactionTemplate: false
        };

        fetchTransactionTemplatesIfNecessary(props);
    }

    componentWillReceiveProps(newProps) {
        fetchTransactionTemplatesIfNecessary(newProps);
    }

    startNewTransactionTemplate = () => {
        this.setState({
            newTransactionTemplate: true
        });
    }

    exitNewTransactionTemplate = () => {
        this.setState({
            newTransactionTemplate: false
        });
    }

    render () {
        const {
            accountId,
            accountTransactionTemplate,
            currency
        } = this.props;

        let transactions = accountTransactionTemplate.get(accountId).get('transactionTemplates');

        return (
            <div>
                <div className={ styles.headings }>
                    <span className={ styles.column }>Template Name</span>
                    <span className={ styles.column }>Name</span>
                    <span className={ styles.column }>Category</span>
                    <span className={ styles.column }>Amount</span>
                </div>
                { !this.state.newTransactionTemplate ?
                        (
                            <button className={ styles.newTransactionTemplate } onClick={ this.startNewTransactionTemplate }>
                                New
                            </button>
                        ) : (
                            <TransactionTemplateForm accountId={ accountId } form='new' done={ this.exitNewTransactionTemplate } currency={ currency } initialValues={ { templateName: '', name: '', category: '', amount: 0 } } />
                        )
                }
                { transactions.map(transaction => (<TransactionTemplate key={ transaction.get('id') } transactionTemplate={ transaction } currency={ currency }/>)).toOrderedSet().toArray() }
            </div>
        );
    }
}
