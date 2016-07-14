import React from 'react';
import { render } from 'react-dom';
import { connect } from 'react-redux';
import { fetchAccounts, importData } from '../../actions';
import AccountList from '../accountList';
import TransactionList from '../transactionList';
import Loader from '../loader';
import styles from './accountsPage.css';

@connect((state) => {
	return {
		accounts: state.accounts,
	}
})
class AccountsPage extends React.Component {
	static propTypes = {
		accounts: React.PropTypes.object.isRequired
	}

	constructor (props) {
		super(props);
		let selected = Object.keys(props.accounts).length !== 0
			? Object.keys(props.accounts)[0]
			: "";
		this.state = {
			selected: selected
		};
	}

	selectAccount = (id) => {
		this.setState({
			selected: id
		});
	}

	importButton = () => {
		this.props.dispatch(importData());
	}

	render () {
		const accounts = this.props.accounts
		const selected = this.state.selected;

		let transactionIds = [];
		let currency;
		if (selected) {
			transactionIds = accounts[selected].transactions;
			currency = accounts[selected].currency;
		}

		return (
			<div className={ styles.accountsPage }>
				{
					Object.keys(accounts).length !== 0 ? (
						<div className={ styles.accountList }>
							<AccountList selected={ selected } onSelect={ this.selectAccount }/>
						</div>
					) : (
						<div>
							Place QIF files in the /import folder and click <button onClick={ this.importButton }>Import</button>
						</div>
					)
				}
				{
					transactionIds.length !== 0 ? (
						<div className={ styles.transactionList }>
							<TransactionList accountId={ selected } transactionIds={ transactionIds } currency={ currency } />
						</div>
					) : null
				}
			</div>
		)
	}
}

@connect((state) => {
	return {
		fetching: state.fetching
	}
})
export default class AccountsPageWrapper extends React.Component {
	static propTypes = {
		fetching: React.PropTypes.bool.isRequired
	}

	constructor (props) {
		super(props);
		props.dispatch(fetchAccounts());
	}

	render () {
		return (
			<Loader loading={ this.props.fetching }>
				<AccountsPage />
			</Loader>
		)
	}
}
