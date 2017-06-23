import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import classNames from 'classnames';
import { fetchSummary } from '../../actions/summary';
import SummaryGraphs from './SummaryGraphs';
import SummaryTransaction from './SummaryTransaction';
import Chart from 'chart.js';
import './Summary.css';

const ONE_WEEK = 7;
const TWO_WEEKS = 14;
const ONE_MONTH = 30;

class Summary extends Component {
    static propTypes = {
        summary: ImmutablePropTypes.list.isRequired,
        dispatch: PropTypes.func.isRequired
    }

    constructor(props) {
        super(props);

        Chart.defaults.global.elements.arc.borderColor = '#312c32';
        Chart.defaults.global.elements.arc.borderWidth = 1;
        this.state = {
            period: TWO_WEEKS
        }
    }

    fetchForPeriod = period => {
        const { dispatch } = this.props;

        const since = new Date();
        since.setDate(since.getDate() - period);
        dispatch(fetchSummary(since));
    }

    componentWillMount() {
        const { period } = this.state;
        this.fetchForPeriod(period);
    }

    selectPeriod = period => () => {
        if (period === this.state.period) {
            return;
        }

        this.setState({ period });
        this.fetchForPeriod(period);
    }

    render() {
        const { period } = this.state;
        const { summary } = this.props;

        return (
                <div className="summary container">
                    <div className="nameAndPeriod">
                        <h1>Summary</h1>
                        <div className="periods">
                            <div className={ classNames('period', { selected: period === ONE_WEEK }) } onClick={ this.selectPeriod(ONE_WEEK) }>
                                1 Week
                            </div>
                            <div className={ classNames('period', { selected: period === TWO_WEEKS }) } onClick={ this.selectPeriod(TWO_WEEKS) }>
                                2 Weeks
                            </div>
                            <div className={ classNames('period', { selected: period === ONE_MONTH }) } onClick={ this.selectPeriod(ONE_MONTH) }>
                                1 Month
                            </div>
                        </div>
                    </div>
                    <div className="divider"></div>
                    {
                    !summary.isEmpty()
                    ? (
                    <div className="graphsAndTransactions">
                        <div className="graphsSection">
                            <h2>Graphs</h2>
                            <SummaryGraphs />
                        </div>
                        <div className="transactionsSection">
                            <h2>Transactions</h2>
                            <div>
                                {
                                summary.map(t => <SummaryTransaction key={ t.get('id') } transaction={ t } />)
                                }
                            </div>
                        </div>
                    </div>
                    )
                    : (
                    <div className="summary container">No transactions in the last
                        {
                        period === ONE_WEEK
                        ? ' one week.'
                        : period === TWO_WEEKS
                        ? ' two weeks.'
                        : ' one month.'
                        }
                    </div>
                    )
                    }
                </div>
                );
    }
}

export default connect(state => ({
    summary: state.summary
}))(Summary);
