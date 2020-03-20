import React from 'react';
import './SignupForm.scss';
import { Field, reduxForm, formValueSelector } from 'redux-form';
import { connect } from 'react-redux'
import { red } from '@material-ui/core/colors';

function SignupForm() {
    return (
        <form>
            <div>
                <label htmlFor="username">Username</label>
                <Field name="username" component="input" type="text" />
            </div>
            <div>
                <label htmlFor="password">Password</label>
                <Field name="password" component="input" type="password" />
            </div>
        </form>
    )
}

export default reduxForm({
    form: 'signup',
})(SignupForm)
