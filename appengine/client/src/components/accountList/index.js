import React from 'react';
import { render } from 'react-dom';
import { connect } from 'react-redux';
import classNames from 'classnames';
import { toCurrency } from '../../utils';
import styles from './accountList.css';

@connect((state) => {
	return {
		accounts: state.accounts
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
					{ Object.keys(accounts).map(id => {
						let account = accounts[id];
						let selectedClass = {};
						selectedClass[styles.selected] = selected === account.id;
						return (
							<button key={ account.id } className={ classNames(styles.accountButton, selectedClass) } onClick={ onSelect.bind(this, id) }>
								<div className={ styles.accountName }>
									{ account.name }
								</div>
								<div className={ styles.accountBalance }>
									Balance: { toCurrency(account.balance, account.currency) }
								</div>
							</button>
						)
					}) }
				</div>
			</div>
		)
	}
}
