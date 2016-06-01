import React from 'react';
import styles from './transaction.css';
import { toCurrency, toDate } from '../../utils';

export default class Transaction extends React.Component {
	static propTypes = {
		transaction: React.PropTypes.object.isRequired,
		currency: React.PropTypes.string.isRequired
	};

	render () {
		const {
			transaction,
			currency
		} = this.props;

		return (
			<div>
				<span className={ styles.transactionField }>{ transaction.name }</span>
				<span className={ styles.transactionField }>{ transaction.type }</span>
				<span className={ styles.transactionField }>{ toDate(transaction.time) }</span>
				<span className={ styles.transactionField }>{ transaction.category }</span>
				<span className={ styles.transactionField }>{ toCurrency(transaction.incoming, currency) }</span>
				<span className={ styles.transactionField }>{ toCurrency(transaction.outgoing, currency) }</span>
			</div>
		)
	}
}

