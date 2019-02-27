/**
 * Copyright 2018 Shift Devices AG
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { Component, h } from 'preact';
import Spinner from '../../../components/spinner/Spinner';
import { apiGet, apiPost } from '../../../utils/request';
import { apiWebsocket } from '../../../utils/websocket';
import { AccountInterface } from '../account';
import { FetchBalances } from './fetchbalances';

interface ProvidedProps {
    accounts: AccountInterface[];
}

interface AccountInitialized {
    [code: string]: boolean;
}

interface State {
    initialized: AccountInitialized;
}

type Props = ProvidedProps;

export class InitializeAllAccounts extends Component<Props, State> {
    constructor(props) {
        super(props);
        const initialized: AccountInitialized = {};
        this.props.accounts.map((account: AccountInterface) => {
            initialized[account.code] = false;
        });
        this.state = { initialized };
      }

    private unsubscribe!: () => void;

    public componentDidMount() {
        this.props.accounts.map((account: AccountInterface) => {
            this.onStatusChanged(account.code);
        });
        this.unsubscribe = apiWebsocket(this.onEvent);
    }

    private allInitialized() {
        return Object.keys(this.state.initialized).every(key => this.state.initialized[key]) ? true : false;
    }

    private onEvent = data => {
        switch (data.data) {
            case 'statusChanged':
                this.onStatusChanged(data.code);
                break;
        }
    }

    private onStatusChanged(code) {
        apiGet(`account/${code}/status`).then(status => {
            const accountSynced = status.includes('accountSynced');
            if (!accountSynced && !status.includes('accountDisabled')) {
                apiPost(`account/${code}/init`);
            }
            this.setState(state => {
                const initialized = state.initialized;
                initialized[code] = accountSynced;
                return { initialized };
            });
        });
    }

    public componentWillUnmount() {
        this.unsubscribe();
    }

    public render(): JSX.Element {
        if (this.allInitialized()) {
            return <FetchBalances accounts={this.props.accounts}/>;
        }
        return (
            <div>
                <Spinner text="Synchronizing all accounts..."/>
            </div>
        );
    }
}
