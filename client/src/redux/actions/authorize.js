export const authorize = data => dispatch => {
    dispatch({
        type: 'AUTHORIZE',
        payload: { authorized: true, admin: data.isAdmin, username: data.username}
    })
}