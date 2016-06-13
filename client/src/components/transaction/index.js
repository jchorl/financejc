import React from 'react';
import classNames from 'classnames';
import { reduxForm } from 'redux-form';
import styles from './transaction.css';
import { toCurrency, toDate } from '../../utils';
import { putTransaction } from '../../actions';

export default class Transaction extends React.Component {
	static propTypes = {
		transaction: React.PropTypes.object.isRequired,
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

		let fields = (
			<div className={ styles.transactionFields }>
				<span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ transaction.name }</span>
				<span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ toDate(transaction.date) }</span>
				<span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ transaction.category }</span>
				<span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ toCurrency(transaction.amount, currency) }</span>
			</div>
		);

		if (this.state.editMode) {
			fields = (
				<TransactionForm form={ transaction.id } transaction={ transaction } done={ this.exitEditMode }/>
			);
		}

		return (
			<div className={ styles.transaction }>
				{ fields }
			</div>
		)
	}
}

function pad(n) {
	return n<10 ? '0'+n : n
}

function toRFC3339(datestring) {
	let d = new Date(datestring)
	return d.getFullYear()+'-'
	+ pad(d.getMonth()+1)+'-'
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
	return {
		initialValues: {
			name: props.transaction.name,
			date: toRFC3339(props.transaction.date),
			category: props.transaction.category,
			amount: props.transaction.amount
		}
	}
})
class TransactionForm extends React.Component {
	static propTypes = {
		fields: React.PropTypes.object.isRequired,
		transaction: React.PropTypes.object.isRequired,
		done: React.PropTypes.func.isRequired
	};

	submit = (data) => {
		let obj = Object.assign({}, this.props.transaction, data);
		obj.date = new Date(obj.date);
		this.props.dispatch(putTransaction(obj));
		this.props.done();
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
