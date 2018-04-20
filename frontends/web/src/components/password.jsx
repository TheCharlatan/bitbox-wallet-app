import { Component } from 'preact';

import { Input, Checkbox, Field } from './forms';

export class PasswordInput extends Component {

    tryPaste = event => {
        if (event.target.type === 'password') {
            event.preventDefault();
            alert('TODO nice message: to paste text, enable \"see plaintext\"');
        }
    }

    render(props) {
        const { seePlaintext, ...rest } = props;
        return (
            <Input
                type={seePlaintext ? 'text' : 'password'}
                autoComplete="off"
                onPaste={this.tryPaste}
                {...rest}
            />
        );
    }
}

export class PasswordRepeatInput extends Component {
    constructor(props) {
        super(props);
        this.state = this.getInitialState();
    }

    clear = () => {
        this.setState(this.getInitialState());
    }

    getInitialState() {
        return {
            password: '',
            passwordRepeat: '',
            seePlaintext: false,
            capsLock: false
        };
    }

    validate = () => {
        if (this.state.password && this.state.password === this.state.passwordRepeat) {
            this.props.onValidPassword(this.state.password);
        } else {
            this.props.onValidPassword(null);
        }
    }

    handleFormChange = event => {
        let value = event.target.value;
        if (event.target.type === 'checkbox') {
            value = event.target.checked;
        }
        this.setState({ [event.target.id]: value });
        this.validate();
    }

    handleCheckCaps = event => {
        const capsLock = event.getModifierState && event.getModifierState('CapsLock');
        this.setState({ capsLock });
    }

    render({ disabled, helptext }, { password, passwordRepeat, seePlaintext, capsLock }) {
        const warning = (capsLock && !seePlaintext) && <p>WARNING: caps lock (⇪) are enabled</p>;

        return (
            <div>
                <PasswordInput
                    autoFocus
                    id="password"
                    seePlaintext={seePlaintext}
                    label="Password"
                    placeholder={helptext}
                    disabled={disabled}
                    onInput={this.handleFormChange}
                    onKeyUp={this.handleCheckCaps}
                    onKeyDown={this.handleCheckCaps}
                    value={password}
                />
                <PasswordInput
                    id="passwordRepeat"
                    seePlaintext={seePlaintext}
                    label="Repeat Password"
                    disabled={disabled}
                    onInput={this.handleFormChange}
                    onKeyUp={this.handleCheckCaps}
                    onKeyDown={this.handleCheckCaps}
                    value={passwordRepeat}
                />
                {warning}
                <Field>
                    <Checkbox
                        id="seePlaintext"
                        onChange={this.handleFormChange}
                        checked={seePlaintext}
                        label="See Plaintext"
                    />
                </Field>
            </div>
        );
    }
}
