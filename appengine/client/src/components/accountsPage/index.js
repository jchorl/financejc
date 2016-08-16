import React from 'react';
import Immutable from 'immutable';
import { render } from 'react-dom';
import { connect } from 'react-redux';
import { fetchAccounts, importData } from '../../actions';
import AccountList from '../accountList';
import TransactionList from '../transactionList';
import Loader from '../loader';
import styles from './accountsPage.css';

@connect((state) => {
	return {
		accounts: state.accountTransaction.get('account')
	}
})
class AccountsPage extends React.Component {
	static propTypes = {
		accounts: React.PropTypes.object.isRequired
	}

	constructor (props) {
		super(props);
		let selected = props.accounts.size !== 0
			? props.accounts.first().get('id')
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

		let transactionIds = Immutable.List();
		let currency;
		if (selected) {
			transactionIds = accounts.get(selected).get('transactions');
			currency = accounts.get(selected).get('currency');
		}

		return (
			<div className={ styles.accountsPage }>
				{
					accounts.size !== 0 ? (
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
					transactionIds.size !== 0 ? (
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
