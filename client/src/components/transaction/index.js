import React from 'react';
import classNames from 'classnames';
import styles from './transaction.css';
import { toCurrency, toDate } from '../../utils';

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
				<span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ toDate(transaction.time) }</span>
				<span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ transaction.category }</span>
				<span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ toCurrency(transaction.amount, currency) }</span>
			</div>
		);

		if (this.state.editMode) {
			fields = (
				<div className={ styles.transactionFields }>
					<input type="text" className={ styles.transactionField } defaultValue={ transaction.name }/>
					<input type="date" className={ styles.transactionField } defaultValue={ toDate(transaction.time) }/>
					<input type="text" className={ styles.transactionField } defaultValue={ transaction.category }/>
					<input type="text" className={ styles.transactionField } defaultValue={ transaction.amount }/>
				</div>
			);
		}

		return (
			<div className={ styles.transaction }>
				{ fields }
				{ this.state.editMode && (
					<div className={ styles.saveExit }>
						<div>
							<button onClick={ this.exitEditMode }>Cancel</button>
							<button onClick={ this.save }>Save</button>
						</div>
					</div>
				)}
			</div>
		)
	}
}

