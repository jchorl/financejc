import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import Chart from 'chart.js';
import { toCurrencyString } from '../../util';

function generateColors(size) {
    let background = [];
    let hover = [];
    for (let i = 0; i < size; ++i) {
        background.push('hsl(' + (360 * i / size) + ', 80%, 70%)');
        hover.push('hsl(' + (360 * i / size) + ', 80%, 60%)');
    }
    return { background, hover };
}

class IncomeSpendGraph extends Component {
    static propTypes = {
        currency: ImmutablePropTypes.contains({
            items: ImmutablePropTypes.mapOf(
                           ImmutablePropTypes.contains({
                               digitsAfterDecimal: PropTypes.number.isRequired
                           }).isRequired
                           ).isRequired
        }).isRequired,
        currencyCode: PropTypes.string.isRequired,
        transactions: ImmutablePropTypes.listOf(
                ImmutablePropTypes.contains({
                    category: PropTypes.string.isRequired,
                    amount: PropTypes.number.isRequired
                }).isRequired
                ).isRequired
    }

    componentWillReceiveProps(nextProps) {
        this.renderChart(this.chartCtx, nextProps);
    }

    renderChart = (chart, props) => {
        if (chart === null) {
            return;
        }
        this.chartCtx = chart;
        if (this.chart) {
            this.chart.destroy();
        }

        const {
            currency,
            currencyCode,
            transactions
        } = props;

        const categoryToIncomeSpendToTotal = transactions
            .groupBy(t => t.get('category'))
            .map(transactionList =>
                    transactionList.groupBy(t => t.get('amount') > 0 ? 'income' : 'spend')
                    .map(transactionList => transactionList.reduce((sum, t) => sum + t.get('amount'), 0))
                    );

        const colors = generateColors(categoryToIncomeSpendToTotal.size);
        const datasets = categoryToIncomeSpendToTotal.mapEntries(([category, incomeSpendToTotal], idx) => [
                category,
                {
                    label: category,
                    data: [ incomeSpendToTotal.get('income', 0), Math.abs(incomeSpendToTotal.get('spend', 0)) ],
                    backgroundColor: [ colors.background[idx], colors.background[idx] ],
                    hoverBackgroundColor: [ colors.hover[idx], colors.hover[idx] ]
                }
        ]).valueSeq().toJS();
        const data = {
            datasets,
            labels: [ 'Income', 'Spend' ]
        };

        const digitsAfterDecimal = currency.getIn(['items', currencyCode, 'digitsAfterDecimal']);
        this.chart = new Chart(chart, {
            type: 'bar',
            data,
            options: {
                legend: { display: false },
                maintainAspectRatio: false,
                scales: {
                    xAxes: [{
                        stacked: true
                    }],
                    yAxes: [{
                        stacked: true
                    }]
                },
                tooltips: {
                    callbacks: {
                        label: (tooltipItem, data) => {
                            let label = data.datasets[tooltipItem.datasetIndex].label + ': ';
                            label += toCurrencyString(tooltipItem.yLabel, currencyCode, digitsAfterDecimal);
                            return label;
                        },
                        title: () => ''
                    }
                }
            }
        });

    }

    render() {
        const { currencyCode } = this.props;
        return (
                <div className="graphWrapper">
                    <div className="graph">
                        <canvas ref={ chart => this.renderChart(chart, this.props) }></canvas>
                    </div>
                    <div className="currencyLabel">{ currencyCode }</div>
                </div>
                );
    }
}

export default connect(state => ({
    currency: state.currency
}))(IncomeSpendGraph);
