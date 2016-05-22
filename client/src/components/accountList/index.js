import React from 'react';
import { render } from 'react-dom';
import styles from './accountList.css';

export default class AccountList extends React.Component {
	static propTypes = {
		accounts: React.PropTypes.array.isRequired
	}

	render () {
		return (
			<div>
				<h3 className={ styles.accountsTitle }>Accounts</h3>
				<div>
					{ this.props.accounts.map(function(account) {
						return (
							<button key={ account.id } className={ styles.accountButton }>
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
