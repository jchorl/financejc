import React, { Component } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import Template from './Template';
import TemplateEdit from './TemplateEdit';

export default class TemplateWrapper extends Component {
    static propTypes = {
        template: ImmutablePropTypes.map.isRequired,
        currency: ImmutablePropTypes.map.isRequired
    }

    constructor() {
        super();
        this.state = { isEditMode: false }
    }

    enterEditMode = () => {
        this.setState({ isEditMode: true });
    }

    exitEditMode = () => {
        this.setState({ isEditMode: false });
    }

    render() {
        const {
            template,
            currency
        } = this.props;

        const { isEditMode } = this.state;

        return isEditMode
            ? (
                    <TemplateEdit template={ template } currency={ currency } exitEditMode={ this.exitEditMode } />
                    )
            : (
                    <Template template={ template } currency={ currency } enterEditMode={ this.enterEditMode } />
                    );
    }
}
