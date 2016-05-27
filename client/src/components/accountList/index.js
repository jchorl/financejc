import React from 'react';
import { render } from 'react-dom';
import classNames from 'classnames';
import styles from './accountList.css';

export default class AccountList extends React.Component {
	static propTypes = {
		accounts: React.PropTypes.array.isRequired,
		onSelect: React.PropTypes.func,
		selected: React.PropTypes.number
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
					{ accounts.map((account, idx) => {
						let selectedClass = {};
						selectedClass[styles.selected] = selected === idx;
						return (
							<button key={ account.id } className={ classNames(styles.accountButton, selectedClass) } onClick={ onSelect.bind(this, idx) }>
								<div className={ styles.accountName }>
									{ account.name }
								</div>
								<div className={ styles.accountBalance }>
									Balance: { account.balance.toLocaleString(window.navigator.language, {style: 'currency', currency: account.currency}) }
								</div>
							</button>
						)
					}) }
				</div>
			</div>
		)
	}
}
