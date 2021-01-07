import { configureStore, getDefaultMiddleware } from '@reduxjs/toolkit';
import reduxLogger from 'redux-logger';
import reduxWebsocket from '@giantmachines/redux-websocket';
import counterReducer from '../features/counter/counterSlice';
import logsReducer from '../features/logs/logsReducer';
import tunnelsReducer from '../features/tunnels/tunnelsReducer';
import authReducer from '../features/auth/authReducer';

export default configureStore({
  reducer: {
    counter: counterReducer,
    logs: logsReducer,
    tunnels: tunnelsReducer,
    auth: authReducer,
  },
  middleware: [
    ...getDefaultMiddleware(),
    reduxLogger,
    reduxWebsocket({
      serializer: JSON.stringify,
      deserializer: JSON.parse
    })
  ]
});
