import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import CategoryGraph from './CategoryGraph';
import IncomeSpendGraph from './IncomeSpendGraph';

class SummaryGraphs extends Component {
    static propTypes = {
        account: ImmutablePropTypes.contains({
            items: ImmutablePropTypes.mapOf(
                           ImmutablePropTypes.contains({
                               currency: PropTypes.string.isRequired
                           }).isRequired
                           ).isRequired
        }).isRequired,
        currency: ImmutablePropTypes.contains({
            items: ImmutablePropTypes.mapOf(
                           ImmutablePropTypes.contains({
                               digitsAfterDecimal: PropTypes.number.isRequired
                           }).isRequired
                           ).isRequired
        }).isRequired,
        summary: ImmutablePropTypes.listOf(
                         ImmutablePropTypes.contains({
                             name: PropTypes.string.isRequired,
                             date: PropTypes.instanceOf(Date),
                             category: PropTypes.string.isRequired,
                             amount: PropTypes.number.isRequired,
                             accountId: PropTypes.number.isRequired
                         }).isRequired
                         ).isRequired
    }

    splitByCurrency = () => {
        const {
            account,
            summary
        } = this.props;

        return summary.groupBy(t => account.getIn(['items', t.get('accountId'), 'currency']));
    }

    render() {
        const splitByCurrency = this.splitByCurrency();

        return (
                <div className="graphSections">
                    <div className="graphSection">
                        <h3>Expenses by Category</h3>
                        <div className="graphs">
                            { splitByCurrency.map((transactions, currency) => <CategoryGraph key={ currency } transactions={ transactions } currencyCode={ currency } />).valueSeq().toJS() }
                        </div>
                    </div>
                    <div className="graphSection">
                        <h3>Income/Spending</h3>
                        <div className="graphs">
                            { splitByCurrency.map((transactions, currency) => <IncomeSpendGraph key={ currency } transactions={ transactions } currencyCode={ currency } />).valueSeq().toJS() }
                        </div>
                    </div>
                </div>
                );
    }
}

export default connect(state => ({
    account: state.account,
    currency: state.currency,
    summary: state.summary
}))(SummaryGraphs);
