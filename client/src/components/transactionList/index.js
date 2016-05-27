import React from 'react';

export default class TransactionList extends React.Component {
	static propTypes = {
		transactions: React.PropTypes.array.isRequired
	};

	static defaultProps = {
		transactions: []
	};

	render () {
		const {
			transactions
		} = this.props;

		return (
			<div>
				{ transactions.map(transaction => (<div>{ transaction.id }</div>)) }
			</div>
		)
	}
}
