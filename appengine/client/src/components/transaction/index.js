import React from 'react';
import classNames from 'classnames';
import { connect } from 'react-redux';
import { reduxForm } from 'redux-form';
import styles from './transaction.css';
import { toCurrency, toDate } from '../../utils';
import { putTransaction } from '../../actions';

@connect((state) => {
	return {
		transactions: state.transactions
	}
})
export class Transaction extends React.Component {
	static propTypes = {
		transactionId: React.PropTypes.string.isRequired,
		transactions: React.PropTypes.object.isRequired,
		currency: React.PropTypes.string.isRequired
	};

	constructor(props) {
		super(props);
		this.state = {editMode: false};
	}

	enterEditMode = () => {
		this.setState({editMode: true});
	}

	exitEditMode = () => {
		this.setState({editMode: false});
	}

	save = () => {
		this.exitEditMode();
	}

	render () {
		const {
			transactions,
			transactionId,
			currency
		} = this.props;

		let transaction = transactions[transactionId];

		return transaction ? (
			<div className={ styles.transaction }>
				{ this.state.editMode ? (
					<TransactionForm form={ transactionId } transaction={ transaction } done={ this.exitEditMode }/>
				) : (
					<div className={ styles.transactionFields }>
						<span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ transaction.name }</span>
						<span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ toDate(transaction.date) }</span>
						<span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ transaction.category }</span>
						<span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ toCurrency(transaction.amount, currency) }</span>
					</div>
				) }
			</div>
		) : null
	}
}

function pad(n) {
	return n<10 ? '0'+n : n
}

function toRFC3339(datestring) {
	let d = new Date(datestring)
	return d.getFullYear() + '-'
	+ pad(d.getMonth() + 1) + '-'
	+ pad(d.getDate());
}

@reduxForm({
	fields: [
		'name',
		'date',
		'category',
		'amount'
	]
},
(state, props) => {
	if (props.transaction) {
		return {
			initialValues: {
				name: props.transaction.name,
				date: toRFC3339(props.transaction.date),
				category: props.transaction.category,
				amount: props.transaction.amount
			}
		};
	}
	return {
		initialValues: {
			name: "",
			date: toRFC3339(new Date()),
			category: "",
			amount: 0
		}
	};
})
export class TransactionForm extends React.Component {
	static propTypes = {
		fields: React.PropTypes.object.isRequired,
		// either transaction (for editing) or accountId (for new transactions) should be passed
		transaction: React.PropTypes.object,
		accountId: React.PropTypes.string,
		done: React.PropTypes.func
	};

	submit = (data) => {
		const {
			accountId,
			dispatch,
			done,
			transaction
		} = this.props;

		console.log(`about to submit accountId: ${accountId}`);

		let obj = Object.assign({}, transaction, data);
		obj.date = new Date(obj.date);
		obj.accountId = accountId;
		obj.amount = parseFloat(obj.amount);
		dispatch(putTransaction(obj));
		done && done();
	}

	render () {
		const {
			fields: {
				name,
				date,
				category,
				amount
			},
			handleSubmit
		} = this.props;

		return (
			<form className={ styles.transactionFields } onSubmit={ handleSubmit(this.submit) }>
				<input type="text" placeholder="Name" className={ styles.transactionField } { ...name }/>
				<input type="date" placeholder={ toRFC3339() } className={ styles.transactionField } { ...date }/>
				<input type="text" placeholder="Category" className={ styles.transactionField } { ...category }/>
				<input type="number" placeholder="0" step="0.01" className={ styles.transactionField } { ...amount }/>
				<div className={ styles.saveExit }>
					<div>
						<button type="button" onClick={ this.props.done }>Cancel</button>
						<button type="submit">Save</button>
					</div>
				</div>
			</form>
		);
	}
}