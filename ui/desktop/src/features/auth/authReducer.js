import { createReducer } from "@reduxjs/toolkit";

export const tunnelsReducer = createReducer(
  { loggedIn: false, authInstructions: null, syncedWithBackend: false },
  {
    "REDUX_WEBSOCKET::MESSAGE": (state, action) => {
      if (action.payload.message.type === "AuthorizationInfo") {
        state.loggedIn = action.payload.message.loggedIn;
        state.syncedWithBackend = true;
      }
      if (action.payload.message.type === "AuthorizationInstructions")
        state.authInstructions = action.payload.message;
    },
    "REDUX_WEBSOCKET::SEND": (state, action) => {
      if (action.payload.messageType === "MT_Logout") {
        state.authInstructions = null;
      }
    },
  }
);

export default tunnelsReducer;
