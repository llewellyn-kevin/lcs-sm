import React from 'react';
import axios from 'axios';
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
            ],
            leagues: [],
            seasons: [],
        };
    }

    componentDidMount() {
        axios.get('http://localhost:8080/api/v2/leagues').then(res => {
            this.setState({ leagues: res.data });
        }).catch(err => {
            console.log(err);
        });

        axios.get('http://localhost:8080/api/v2/seasons').then(res => {
            const seasons = res.data.map(s => {
                s.name = s.year + ' ' + s.season;
                return s;
            });
            this.setState({ seasons });
        }).catch(err => {
            console.log('error', err);
        });
    }

    render() {
        const listItems = this.state.teams.map((team) =>
            <StockTickerItem key={team.id.toString()} team={team} />
        );

        return(
            <aside className="StockTicker">
                <div className="StockTicker-CustomSelect-row">
                    <CustomSelect
                        label="LEAGUE"
                        options={this.state.leagues} />
                    <CustomSelect
                        label="SEASON"
                        options={this.state.seasons} />
                </div>
                {listItems}
            </aside>
        );
    }
}