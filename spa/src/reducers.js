import { combineReducers } from 'redux';

const topPage = (state = {}, action) => {
  switch (action.type) {
  case 'SET_TOP_PAGE_MESSAGE':
    return Object.assign({}, state, {
      message: action.message
    })
  default:
    return state
  }
}

export const spaApp = combineReducers({
  topPage
})
