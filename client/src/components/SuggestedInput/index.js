import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { List } from 'immutable';
import { connect } from 'react-redux';
import { fetchSuggestions } from '../../actions/autocomplete';
import './SuggestedInput.css';

class SuggestedInput extends Component {
    static propTypes = {
        autocomplete: ImmutablePropTypes.contains({
            field: PropTypes.string.isRequired,
            term: PropTypes.string.isRequired,
            items: ImmutablePropTypes.list.isRequired
        }),
        accountId: PropTypes.number.isRequired,
        dispatch: PropTypes.func.isRequired,
        field: PropTypes.string.isRequired,
        onChange: PropTypes.func,
        selectSuggestion: PropTypes.func.isRequired,
        hoverSuggestion: PropTypes.func,
        unhoverSuggestion: PropTypes.func
    }

    constructor(props) {
        super(props);

        this.propsToFilterOut = [
            'accountId',
            'autocomplete',
            'dispatch',
            'field',
            'onChange',
            'selectSuggestion',
            'hoverSuggestion',
            'unhoverSuggestion'
        ];
        this.state = {
            suggestions: List(),
            nextField: '',
            nextTerm: '',
            alreadySelected: false,
            hasHovered: false
        };
    }

    componentWillReceiveProps(nextProps) {
        const {
            nextField,
            nextTerm
        } = this.state;

        const { autocomplete } = nextProps;

        if (this.props.autocomplete !== autocomplete &&
                autocomplete.get('field') === nextField &&
                autocomplete.get('term') === nextTerm) {
            this.setState({
                suggestions: autocomplete.get('items'),
                alreadySelected: false
            });
        }
    }

    onChange = e => {
        const {
            accountId,
            dispatch,
            field,
            onChange
        } = this.props;

        this.setState({
            nextField: field,
            nextTerm: e.target.value
        });

        dispatch(fetchSuggestions(accountId, field, e.target.value));
        if (onChange) {
            onChange(e);
        }
    }

    selectSuggestion = s => () => {
        // clear out suggestions to hide box
        this.setState({
            suggestions: List(),
            alreadySelected: true,
            hasHovered: false
        });

        this.props.selectSuggestion(s);
    }

    hoverSuggestion = s => () => {
        const { hoverSuggestion } = this.props;

        this.setState({ hasHovered: true });
        if (hoverSuggestion) {
            hoverSuggestion(s);
        }
    }

    unhoverSuggestion = () => {
        const { unhoverSuggestion } = this.props;

        this.setState({ hasHovered: false });
        if (unhoverSuggestion) {
            unhoverSuggestion();
        }
    }

    onBlur = () => {
        const { alreadySelected, hasHovered } = this.state;
        // onBlur gets called after the user picks a selection
        if (alreadySelected) {
            return;
        }

        const { unhoverSuggestion } = this.props;

        // only unhover if the user has actually hovered
        if (unhoverSuggestion && hasHovered) {
            unhoverSuggestion();
        }

        // clear out suggestions to hide box
        this.setState({
            suggestions: List(),
            hasHovered: false
        });
    }

    render() {
        let filteredProps = Object.assign({}, this.props);
        for (let propToFilter of this.propsToFilterOut) {
            delete filteredProps[propToFilter];
        }

        const { field } = this.props;
        const { suggestions } = this.state;
        return (
                <div className="suggestedInput">
                    <input onChange={ this.onChange } onBlur={ this.onBlur } { ...filteredProps }></input>
                    {
                    !suggestions.isEmpty()
                    ? (
                        <div className="suggestionBoxWrapper">
                            <div className="suggestionBox" onMouseLeave={ this.unhoverSuggestion }>
                                { suggestions.map(s =>
                                <div
                                    key={ s.get('id') }
                                    className="suggestion"
                                    onMouseDown={ this.selectSuggestion(s) }
                                    onMouseEnter={ this.hoverSuggestion(s) }
                                >{ s.get(field) }</div>
                                ) }
                            </div>
                        </div>
                    )
                    : null
                    }
                </div>
                );
    }
}

export default connect(state => {
    return {
        autocomplete: state.autocomplete
    };
})(SuggestedInput);
