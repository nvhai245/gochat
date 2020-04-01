import { createStore, applyMiddleware, compose } from 'redux';
import thunk from 'redux-thunk';
import rootReducer from './reducers/rootReducer';
import { persistStore, persistReducer } from 'redux-persist';
import storage from 'redux-persist/lib/storage';

const persistConfig = {
    key: 'root',
    storage,
    blacklist: ['form']
}

const persistedReducer = persistReducer(persistConfig, rootReducer);

const state = { 
    authorization: { 
        authorized: false, 
        admin: false, 
        username: "",
    }
}
const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;
export const store = createStore(
        persistedReducer,
        state,
        composeEnhancers(applyMiddleware(thunk))
    );
export const persistor = persistStore(store);