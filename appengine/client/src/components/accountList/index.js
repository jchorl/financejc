import React from 'react';
import { render } from 'react-dom';
import { connect } from 'react-redux';
import classNames from 'classnames';
import { toCurrency } from '../../utils';
import styles from './accountList.css';

@connect((state) => {
	return {
		accounts: state.accountTransaction.get('account')
	}
})
export default class AccountList extends React.Component {
	static propTypes = {
		accounts: React.PropTypes.object.isRequired,
		onSelect: React.PropTypes.func,
		selected: React.PropTypes.string
	}

	render () {
		const {
			accounts,
			onSelect,
			selected
		} = this.props;

		return (
			<div>
				<h3 className={ styles.accountsTitle }>Accounts</h3>
				<div>
					{ accounts.map(account => {
						let selectedClass = {};
						selectedClass[styles.selected] = selected === account.get('id');
						return (
							<button key={ account.get('id') } className={ classNames(styles.accountButton, selectedClass) } onClick={ onSelect.bind(this, account.get('id')) }>
								<div className={ styles.accountName }>
									{ account.get('name') }
								</div>
								<div className={ styles.accountBalance }>
									Balance: { toCurrency(account.get('balance'), account.get('currency')) }
								</div>
							</button>
						)
					}).valueSeq().toArray() }
				</div>
			</div>
		)
	}
}
