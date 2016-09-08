import React from 'react';
import { connect } from 'react-redux';
import styles from './transactionList.css';
import { Transaction, TransactionForm } from '../transaction';

@connect((state) => {
	return {
		accountTransaction: state.accountTransaction
	}
})
export default class TransactionList extends React.Component {
	static propTypes = {
		accountId: React.PropTypes.string.isRequired,
		accountTransaction: React.PropTypes.object.isRequired,
		currency: React.PropTypes.string.isRequired
	};

	constructor (props) {
		super(props);
		this.state = {
			newTransaction: false
		};
	}

	startNewTransaction = () => {
		this.setState({
			newTransaction: true
		});
	}

	exitNewTransaction = () => {
		this.setState({
			newTransaction: false
		});
	}

	render () {
		const {
			accountId,
			accountTransaction,
			currency
		} = this.props;

		let transactions = accountTransaction.get(accountId).get("transactions");

		return (
			<div>
				<div className={ styles.headings }>
					<span className={ styles.column }>Name</span>
					<span className={ styles.column }>Date</span>
					<span className={ styles.column }>Category</span>
					<span className={ styles.column }>Amount</span>
				</div>
				{ !this.state.newTransaction ?
					(
						<button className={ styles.newTransaction } onClick={ this.startNewTransaction }>
							New
						</button>
					) : (
						<TransactionForm accountId={ accountId } form='new' done={ this.exitNewTransaction }/>
					)
				}
				<div>
					{ transactions.map(transaction => (<Transaction key={ transaction.get('id') } transaction={ transaction } currency={ currency }/>)).toArray() }
				</div>
			</div>
		)
	}
}
