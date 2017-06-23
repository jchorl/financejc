import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import PropTypes from 'prop-types';
import classNames from 'classnames';

class CreateAccountEntry extends Component {
    static propTypes = {
        history: PropTypes.shape({
            push: PropTypes.func.isRequired
        }),
        selected: PropTypes.bool.isRequired
    }

    goToCreateAccount = () => {
        const {
            history
        } = this.props;

        history.push('/accounts/create');
    }

    render() {
        const { selected } = this.props;

        return (
                <div className="entry iconLabelEntry" onClick={ this.goToCreateAccount } >
                    <div className={ classNames("selectedBlock", { selected }) }></div>
                    <div className="label"><i className="fa fa-plus-circle"></i><span className="text">New Account</span></div>
                </div>
                );
    }
}

export default withRouter(CreateAccountEntry);
