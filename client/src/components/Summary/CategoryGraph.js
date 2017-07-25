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
        background.push('hsl(' + (360 * i / size) + ', 80%, 80%)');
        hover.push('hsl(' + (360 * i / size) + ', 80%, 60%)');
    }
    return { background, hover };
}

class CategoryGraph extends Component {
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

        const totalsByCategory = transactions.groupBy(t => t.get('category'))
            .map(transactionList => transactionList.reduce(
                        (prevSum, t) => t.get('amount') + prevSum,
                        0
                        ));
        // only want net negative categories
        const totalsByCategoryFiltered = totalsByCategory.filter(v => v < 0).map(v => Math.abs(v));
        const colors = generateColors(totalsByCategoryFiltered.size);
        const data = {
            datasets: [{
                backgroundColor: colors.background,
                hoverBackgroundColor: colors.hover,
                data: totalsByCategoryFiltered.valueSeq().toJS()
            }],
            labels: totalsByCategoryFiltered.keySeq().toJS()
        };

        const digitsAfterDecimal = currency.getIn(['items', currencyCode, 'digitsAfterDecimal']);
        this.chart = new Chart(chart, {
            type: 'pie',
            data,
            options: {
                legend: { display: false },
                tooltips: {
                    callbacks: {
                        label: (tooltipItem, data) => {
                            let dataLabel = data.labels[tooltipItem.index] + ': ';
                            let value = data.datasets[tooltipItem.datasetIndex].data[tooltipItem.index];
                            value = toCurrencyString(value, currencyCode, digitsAfterDecimal);
                            dataLabel += value;

                            return dataLabel;
                        }
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
}))(CategoryGraph);
