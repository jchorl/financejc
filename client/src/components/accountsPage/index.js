import React from 'react';
import { render } from 'react-dom';
import { connect } from 'react-redux';
import { fetchAccounts } from '../../actions';
import AccountList from '../accountList';
import Loader from '../loader';
import styles from './accountsPage.css';

@connect((state) => {
	return { accounts: state.accounts }
})
export default class AccountsPage extends React.Component {
	constructor(props) {
		super(props);
		props.dispatch(fetchAccounts());
	}

	render () {
		return (
			<Loader loading={ this.props.accounts.isFetching }>
				<div className={ styles.accountList }>
					<AccountList accounts={ this.props.accounts.items } />
				</div>
			</Loader>
		)
	}
}
