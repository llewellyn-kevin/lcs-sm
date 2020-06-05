import React from 'react';

// CustomSelect is a custom implementation of a basic html select
// box with unique formatting and properties.
export class CustomSelect extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            value: 0
        };

        this.incrSelection = this.incrSelection.bind(this);
    }

    // temporarily increment selection on click instead of having
    // options appear. Notify parent of change.
    incrSelection(e) {
        let newValue = this.state.value + 1;
        if(newValue >= this.props.options.length) {
            newValue = 0;
        }
        this.setState({ value: newValue });
        this.props.onValueChange(newValue);
    }

    render() {
        return (this.props.options.length > 0)
        ? (
            <div onClick={this.incrSelection} className="CustomSelect btn btn-primary">
                <div className="CustomSelect-value">{this.props.options[this.state.value].name}</div>
                <div className="CustomSelect-arrow">⏷</div>
                <div className="CustomSelect-label">{this.props.label}</div>
            </div>
        )
        : (
            <div className="CustomSelect btn btn-primary">
                <div className="CustomSelect-value">fetching...</div>
            </div>
        );
    }
}