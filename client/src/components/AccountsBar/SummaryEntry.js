import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import PropTypes from 'prop-types';
import classNames from 'classnames';

class SummaryEntry extends Component {
    static propTypes = {
        history: PropTypes.shape({
            push: PropTypes.func.isRequired
        }),
        selected: PropTypes.bool.isRequired
    }

    goToSummary = () => {
        const {
            history
        } = this.props;

        history.push('/accounts/summary');
    }

    render() {
        const { selected } = this.props;

        return (
                <div className="entry iconLabelEntry" onClick={ this.goToSummary } >
                    <div className={ classNames("selectedBlock", { selected }) }></div>
                    <div className="label"><i className="fa fa-pie-chart"></i><span className="text">Summary</span></div>
                </div>
                );
    }
}

export default withRouter(SummaryEntry);
