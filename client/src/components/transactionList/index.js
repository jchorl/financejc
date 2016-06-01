import React from 'react';
import Transaction from '../transaction';
import styles from './transactionList.css';

export default class TransactionList extends React.Component {
	static propTypes = {
		transactions: React.PropTypes.array.isRequired,
		currency: React.PropTypes.string.isRequired
	};

	static defaultProps = {
		transactions: [],
		currency: 'USD'
	};

	render () {
		const {
			transactions,
			currency
		} = this.props;

		return (
			<div>
				<div className={ styles.headings }>
					<span className={ styles.column }>Name</span>
					<span className={ styles.column }>Type</span>
					<span className={ styles.column }>Date</span>
					<span className={ styles.column }>Category</span>
					<span className={ styles.column }>Incoming</span>
					<span className={ styles.column }>Outgoing</span>
				</div>
				<div>
					{ transactions.map(transaction => (<Transaction transaction={ transaction } currency={ currency }/>)) }
				</div>
			</div>
		)
	}
}
