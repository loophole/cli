import { configureStore, getDefaultMiddleware } from "@reduxjs/toolkit";
import reduxLogger from "redux-logger";
import reduxWebsocket from "@giantmachines/redux-websocket";
import logsReducer from "../features/logs/logsReducer";
import tunnelsReducer from "../features/tunnels/tunnelsReducer";
import configReducer from "../features/config/configReducer";

const reducers = {
  logs: logsReducer,
  tunnels: tunnelsReducer,
  config: configReducer,
};

const store = configureStore({
  reducer: reducers,
  middleware: [
    ...getDefaultMiddleware(),
    reduxLogger,
    reduxWebsocket({
      serializer: JSON.stringify,
      deserializer: JSON.parse,
    }),
  ],
});

export default store;
