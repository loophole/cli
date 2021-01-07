import { createReducer } from "@reduxjs/toolkit";

export const logsReducer = createReducer(
  { tunnelLogs: [], globalLogs: [] },
  {
    "REDUX_WEBSOCKET::MESSAGE": (state, action) => {
      if (action.payload.message.type === "Log") {
        state.tunnelLogs.push(action.payload.message);
        state.globalLogs.push(action.payload.message);
      }
      if (action.payload.message.type === "TunnelShutDown") state.tunnelLogs = [];
    },
  }
);

export default logsReducer;
