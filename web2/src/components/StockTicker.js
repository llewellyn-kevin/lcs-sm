import React from 'react';
import {CustomSelect} from './CustomSelect.js';
import {StockTickerItem} from './StockTickerItem.js';

// StockTicker displays StockTickerItems for each of the currently
// selected League and Season.
export class StockTicker extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            teams: [
                {
                    id: 0,
                    name: "Cloud9",
                    code: "C9",
                    values: [758, 721]
                },
                {
                    id: 1,
                    name: "Team Liquid",
                    code: "TL",
                    values: [566, 590]
                },
                {
                    id: 2,
                    name: "Counter Logic Gaming",
                    code: "CLG",
                    values: [540],
                }
            ]
        };
    }

    render() {
        const listItems = this.state.teams.map((team) =>
            <StockTickerItem key={team.id.toString()} team={team} />
        );

        return(
            <aside className="StockTicker">
                <CustomSelect
                    label="LEAGUE"
                    options={["LCS", "LEC", "LCK", "LPL"]} />
                {listItems}
            </aside>
        );
    }
}