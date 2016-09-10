import React from 'react';
import classNames from 'classnames';
import { reduxForm } from 'redux-form';
import styles from './transaction.css';
import { toCurrency, toDate } from '../../utils';
import { putTransaction } from '../../actions';

export class Transaction extends React.Component {
	static propTypes = {
		transaction: React.PropTypes.object,
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
			transaction,
			currency
		} = this.props;

		return transaction ? (
			<div className={ styles.transaction }>
				{ this.state.editMode ? (
					<TransactionForm form={ transaction.get('id') } transaction={ transaction } done={ this.exitEditMode }/>
				) : (
					<div className={ styles.transactionFields }>
						<span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ transaction.get('name') }</span>
						<span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ toDate(transaction.get('date')) }</span>
						<span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ transaction.get('category') }</span>
						<span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ toCurrency(transaction.get('amount'), currency) }</span>
					</div>
				) }
			</div>
		) : null
	}
}

function pad(n) {
	return n<10 ? '0'+n : n
}

function toRFC3339(d) {
	return d.getUTCFullYear() + '-'
	+ pad(d.getUTCMonth() + 1) + '-'
	+ pad(d.getUTCDate());
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
				name: props.transaction.get('name'),
				date: toRFC3339(props.transaction.get('date')),
				category: props.transaction.get('category'),
				amount: props.transaction.get('amount')
			}
		};
	}
	return {
		initialValues: {
			name: '',
			date: toRFC3339(new Date()),
			category: '',
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

		let obj = data;
		if (transaction) {
			obj = Object.assign(transaction.toObject(), data);
		}
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
				<input type="date" placeholder={ toRFC3339(new Date()) } className={ styles.transactionField } { ...date }/>
				<input type="text" placeholder="Category" className={ styles.transactionField } { ...category }/>
				<input type="text" placeholder="0" className={ styles.transactionField } { ...amount }/>
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
