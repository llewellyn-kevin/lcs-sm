import React from 'react';

export class CustomSelect extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            value: 0
        };

        this.incrSelection = this.incrSelection.bind(this);
    }

    incrSelection(e) {
        let newValue = this.state.value + 1;
        if(newValue >= this.props.options.length) {
            newValue = 0;
        }
        this.setState({ value: newValue });
    }

    render() {
        return(
            <div className="CustomSelect btn btn-primary">
                <div className="CustomSelect-value">{this.props.options[this.state.value]}</div>
                <div className="CustomSelect-arrow" onClick={this.incrSelection}>⏷</div>
                <div className="CustomSelect-label">{this.props.label}</div>
            </div>
        );
    }
}