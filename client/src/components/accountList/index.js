import React from 'react';
import { render } from 'react-dom';

export default class AccountList extends React.Component {
	static propTypes = {
		accounts: React.PropTypes.array.isRequired
	}

	render () {
		return (
			<div>
				<h3>Accounts</h3>
				{ this.props.accounts.map(function(account) {
					return <button key={ account.id }>{ account.name }</button>
				}) }
			</div>
		)
	}
}
