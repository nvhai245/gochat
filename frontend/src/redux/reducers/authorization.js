export default (state = {}, action) => {
    switch (action.type) {
        case 'AUTHORIZE':
            return action.payload
        default:
            return state
    }
}