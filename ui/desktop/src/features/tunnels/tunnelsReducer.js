import { createReducer } from "@reduxjs/toolkit";

export const tunnelsReducer = createReducer(
  { tunnel: null },
  {
    "REDUX_WEBSOCKET::MESSAGE": (state, action) => {
      if (action.payload.message.type === "TunnelMetadata") state.tunnel = action.payload.message;
      if (action.payload.message.type === "TunnelShutDown") state.tunnel = null;
    },
  }
);

export default tunnelsReducer;
