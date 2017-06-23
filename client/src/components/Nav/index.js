import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { withRouter } from 'react-router';
import classNames from 'classnames';
import { clearSearchResults, search } from '../../actions/search';
import { UNAUTHD_ID } from '../../constants';
import './Nav.css';

class Nav extends Component {
    static propTypes = {
        user: ImmutablePropTypes.contains({
            id: PropTypes.number.isRequired,
            email: PropTypes.string.isRequired
        }),
        dispatch: PropTypes.func.isRequired,
        history: PropTypes.shape({
            push: PropTypes.func.isRequired
        }),
        location: PropTypes.shape({
            pathname: PropTypes.string.isRequired
        }).isRequired
    }

    constructor(props) {
        super(props);

        this.state = {
            searchTerm: '',
            prevPath: ''
        }
    }

    searchChange = e => {
        const {
            location,
            history
        } = this.props;

        this.props.dispatch(search(e.target.value));
        this.setState({ searchTerm: e.target.value });

        if (location.pathname !== '/search') {
            this.setState({
                prevPath: location.pathname
            });
            history.push('/search');
        }
    }

    returnToPage = () => {
        const {
            dispatch,
            history
        } = this.props;

        const prevPath = this.state.prevPath;
        this.setState({
            prevPath: '',
            searchTerm: ''
        });
        dispatch(clearSearchResults());
        history.push(prevPath);
    }

    render() {
        const {
            location,
            user
        } = this.props;

        const isOnSearchPage = location.pathname === '/search';

        return (
                <div className="nav">
                    <div className="logo">
                        FinanceJC
                    </div>
                    {
                    user.get('id') !== UNAUTHD_ID
                    ? (
                    <div className="searchWrapper">
                        <div className="searchInputIcons">
                            <input type="text" className="searchInput" name="search" placeholder="Search..." onChange={ this.searchChange } value={ this.state.searchTerm } />
                            <i className={ classNames('fa fa-times-circle fa-lg closeIcon', { closeIconShow: isOnSearchPage }) } onClick={ this.returnToPage }></i>
                            <i className="fa fa-search searchIcon"></i>
                        </div>
                    </div>
                    )
                    : null
                    }
                    <div>
                        {
                        user.get('id') !== UNAUTHD_ID
                        ? user.get('email')
                        : ''
                        }
                    </div>
                </div>
                );
    }
}

export default withRouter(connect(state => ({
    user: state.user
}))(Nav));
