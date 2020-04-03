import React from 'react';
import './SignupForm.scss';
import { Field, reduxForm, formValueSelector } from 'redux-form';
import { connect } from 'react-redux'
import { red } from '@material-ui/core/colors';


function SignupForm() {
}

export default reduxForm({
    form: 'signup',
})(SignupForm)
