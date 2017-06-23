import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { Map } from 'immutable';
import { NEW_TEMPLATE_ID } from '../../constants';
import TemplateWrapper from './TemplateWrapper';
import { fetchTemplates } from '../../actions/accountData';
import TemplateEdit from './TemplateEdit.js';
import './TemplateList.css';

function emptyTemplate(accountId) {
    return Map({
        id: NEW_TEMPLATE_ID,
        templateName: '',
        name: '',
        category: '',
        amount: 0,
        accountId
    });
}

class TemplateList extends Component {
    static propTypes = {
        account: ImmutablePropTypes.contains({
            items: ImmutablePropTypes.mapOf(
                           ImmutablePropTypes.contains({
                               id: PropTypes.number.isRequired,
                               currency: PropTypes.string.isRequired
                           }).isRequired
                           ).isRequired
        }).isRequired,
        accountData: ImmutablePropTypes.mapOf(
                             ImmutablePropTypes.contains({
                                 templates: ImmutablePropTypes.contains({
                                     fetched: PropTypes.bool.isRequired,
                                     fetching: PropTypes.bool.isRequired,
                                     items: ImmutablePropTypes.list
                                 }).isRequired
                             }).isRequired
                             ).isRequired,
        currency: ImmutablePropTypes.contains({
            items: ImmutablePropTypes.mapOf(
                    ImmutablePropTypes.contains({
                        digitsAfterDecimal: PropTypes.number.isRequired
                    })
                    ).isRequired
        }).isRequired
    }

    constructor(props) {
        super(props);
        this.state = {
            isEnteringTemplate: false
        };
    }

    fetchTemplatesIfNotFetched = (accountId, props) => {
        const {
            accountData,
            dispatch,
        } = props;

        if (!accountData.getIn([accountId, 'templates', 'fetched']) && !accountData.getIn([accountId, 'templates', 'fetching'])) {
            dispatch(fetchTemplates(accountId));
        }
    }

    componentWillMount() {
        const {
            match: { params: { id } }
        } = this.props;

        let accountId = parseInt(id, 10);
        this.fetchTemplatesIfNotFetched(accountId, this.props);
    }

    componentWillReceiveProps(nextProps) {
        const {
            match: { params: { id } }
        } = nextProps;

        let accountId = parseInt(id, 10);
        this.fetchTemplatesIfNotFetched(accountId, nextProps);
    }

    newTemplate = () => {
        this.setState({
            isEnteringTemplate: true
        });
    }

    doneEnteringTemplate = () => {
        this.setState({
            isEnteringTemplate: false
        });
    }

    render() {
        const {
            account,
            accountData,
            currency,
            match: { params: { id } }
        } = this.props;
        const { isEnteringTemplate } = this.state;

        const parsedId = parseInt(id, 10);
        const accountCurrencyCode = account.getIn(['items', parsedId, 'currency']);
        const accountCurrency = currency.getIn(['items', accountCurrencyCode]);

        return (
                <div className="templateList">
                    {
                    isEnteringTemplate
                    ? (
                    <TemplateEdit
                        template={ emptyTemplate(parsedId) }
                        currency={ accountCurrency }
                        exitEditMode={ this.doneEnteringTemplate }
                    />
                    )
                    : (
                    <div className="newTemplateBar">
                        <button key="EMPTY" onClick={ this.newTemplate }><i className="fa fa-plus-circle"></i> New Template</button>
                    </div>
                    )
                    }
                    {
                    accountData.getIn([parsedId, 'templates', 'items']).map(
                    t => <TemplateWrapper key={ t.get('id') } template={ t } currency={ accountCurrency } />
                    )
                    }
                </div>
                );
    }
}

export default connect(state => {
    return {
        account: state.account,
        accountData: state.accountData,
        currency: state.currency
    }
})(TemplateList);
