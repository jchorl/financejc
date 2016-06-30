import React from 'react';
import Transaction from '../transaction';
import styles from './transactionList.css';

export default class TransactionList extends React.Component {
	static propTypes = {
		transactionIds: React.PropTypes.array.isRequired,
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

	render () {
		const {
			transactionIds,
			transactions,
			currency
		} = this.props;

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
						<div>
							New Transaction Form
						</div>
					)
				}
				<div>
					{ transactionIds.map(transactionId => (<Transaction key={ transactionId } transactionId={ transactionId } currency={ currency }/>)) }
				</div>
			</div>
		)
	}
}
