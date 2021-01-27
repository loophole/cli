import { createReducer } from "@reduxjs/toolkit";
import jwtDecode from "jwt-decode";
import {
  MessageTypeAppStart,
  MessageTypeLogin,
  MessageTypeLoginSuccess,
  MessageTypeLogoutSuccess,
  MessageTypeRequestLogin,
  MessageTypeRequestLogout,
} from "../../constants/websocket";

const defaultDisplayConfig = {
  qr: false,
  verbose: false,
};

export const configReducer = createReducer(
  {
    loggedIn: false,
    authInstructions: null,
    syncedWithBackend: false,
    displayConfig: defaultDisplayConfig,
    feedbackFormUrl: "https://bit.ly/3mvmZBA",
    version: "development",
    commitHash: "unknown",
    homeDirectory: "",
  },
  {
    "REDUX_WEBSOCKET::MESSAGE": (state, action) => {
      if (action.payload.message.type === MessageTypeAppStart) {
        state.loggedIn = action.payload.message.loggedIn;
        state.displayConfig = action.payload.message.displayConfig;
        if (action.payload.message.feedbackFormUrl)
          state.feedbackFormUrl = action.payload.message.feedbackFormUrl;
        state.version = action.payload.message.version;
        state.commitHash = action.payload.message.commitHash;
        state.homeDirectory = action.payload.message.homeDirectory;

        state.user = action.payload.message.idToken
          ? jwtDecode(action.payload.message.idToken)
          : null;

        state.syncedWithBackend = true;
      }
      else if (action.payload.message.type === MessageTypeLogin)
        state.authInstructions = action.payload.message;
      else if (action.payload.message.type === MessageTypeLoginSuccess) {
        state.authInstructions = null;
        state.loggedIn = true;
        state.user = action.payload.message.idToken
          ? jwtDecode(action.payload.message.idToken)
          : null;
      } else if (action.payload.message.type === MessageTypeLogoutSuccess) {
        state.loggedIn = false;
        state.user = null;
        state.syncedWithBackend = true;
      }
    },
    "REDUX_WEBSOCKET::SEND": (state, action) => {
      if (action.payload.type === MessageTypeRequestLogin) {
        state.authInstructions = null;
      } else if (action.payload.type === MessageTypeRequestLogout) {
        state.loggedIn = false;
        state.syncedWithBackend = false;
      } 
    },
  }
);

export default configReducer;
