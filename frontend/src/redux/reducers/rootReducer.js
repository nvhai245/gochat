import { combineReducers } from 'redux';
import authorization from './authorization';
import { reducer as formReducer } from 'redux-form';
export default combineReducers({
 authorization,
 form: formReducer,
});