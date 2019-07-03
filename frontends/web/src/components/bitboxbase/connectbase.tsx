/**
 * Copyright 2019 Shift Devices AG
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

import { Component, h, RenderableProps } from 'preact';
import { alertUser } from '../../components/alert/Alert';
import { Button } from '../../components/forms';
import { apiSubscribe } from '../../utils/event';
import { apiPost } from '../../utils/request';

export interface ConnectedBaseProps {
    bitboxBaseID: string;
}

interface MiddlewareInfoType {
    blocks: number;
    difficulty: number;
    lightningAlias: string;
}

interface State {
    middlewareInfo?: MiddlewareInfoType;
    bitboxBaseID: string;
}

type Props = ConnectedBaseProps;

export class ConnectedBase extends Component<Props, State> {

    constructor(props) {
        super(props);
        this.state = {
            middlewareInfo: undefined,
            bitboxBaseID: '',
        };
    }

    public componentDidMount() {
        // Only create a new websocket if the bitboxBaseID changed.
        if (this.props.bitboxBaseID !== this.state.bitboxBaseID) {
            this.setState({ bitboxBaseID : this.props.bitboxBaseID});
            apiSubscribe('/bitboxbases/' + this.props.bitboxBaseID + '/middlewareinfo', ({ object }) => {
                this.setState({ middlewareInfo: object });
            });
        }
    }

    private removeBitBoxBase = () => {
        apiPost('bitboxbases/disconnectbase', {
            bitboxBaseID : this.props.bitboxBaseID,
        }).then(({ success }) => {
            if (!success) {
                alertUser('Did not work');
            }
        });

    }

    private connectElectrum = () => {
        apiPost('bitboxbases/' + this.props.bitboxBaseID + '/connect-electrum', {
            bitboxBaseID : this.props.bitboxBaseID,
        }).then(({success}) => {
            if (!success) {
                alertUser(success.errorMessage);
            }
        });
    }

    public render(
        {
            bitboxBaseID,
        }: RenderableProps<Props>,
        {
            middlewareInfo,
        }: State,
    ) {
        if (!middlewareInfo) {
            return null;
        }

        return (
                <div class="row">
                    <div class="flex flex-1 flex-row flex-between flex-items-center spaced">
                        <p>Block Number: {middlewareInfo.blocks}</p>
                        <p>Difficulty: {middlewareInfo.difficulty}</p>
                        <p>Device ID: {bitboxBaseID}</p>
                        <p>Lightning Alias: {middlewareInfo.lightningAlias}</p>
                        <div class="buttons flex flex-row flex-end">
                            <Button onClick={this.removeBitBoxBase} danger>Delete</Button>
                        </div>
                    </div>
                    <div class="row">
                        <div class="buttons flex flex-row flex-end">
                            <Button onClick={this.connectElectrum}>Connect Electrum</Button>
                        </div>
                    </div>
                </div>
            );
    }
}
