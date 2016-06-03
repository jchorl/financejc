import React from 'react';
import { render } from 'react-dom';
import { connect } from 'react-redux';
import { fetchAccounts } from '../../actions';
import AccountList from '../accountList';
import TransactionList from '../transactionList';
import Loader from '../loader';
import styles from './accountsPage.css';

class AccountsPage extends React.Component {
	static propTypes = {
		accounts: React.PropTypes.array.isRequired
	}

	constructor (props) {
		super(props);
		this.state = {
			selected: 0
		};
	}

	selectAccount = (index) => {
		this.setState({
			selected: index
		});
	}

	render () {
		const {
			accounts
		} = this.props;

		let transactions = [];
		let currency;
		if (this.state.selected < accounts.length) {
			transactions = accounts[this.state.selected].transactions;
			currency = accounts[this.state.selected].currency;
		}

		return (
			<div className={ styles.accountsPage }>
				<div className={ styles.accountList }>
					<AccountList accounts={ accounts } selected={ this.state.selected } onSelect={ this.selectAccount }/>
				</div>
				<div className={ styles.transactionList }>
					<TransactionList transactions={ transactions } currency={ currency } />
				</div>
			</div>
		)
	}
}

@connect((state) => {
	return { accounts: state.accounts }
})
export default class AccountsPageWrapper extends React.Component {
	static propTypes = {
		accounts: React.PropTypes.object.isRequired
	}

	constructor (props) {
		super(props);
		props.dispatch(fetchAccounts());
	}

	render () {
		return (
			<Loader loading={ this.props.accounts.isFetching }>
				<AccountsPage accounts={ this.props.accounts.items } />
			</Loader>
		)
	}
}
