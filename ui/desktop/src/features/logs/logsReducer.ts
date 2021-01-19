import { createReducer } from "@reduxjs/toolkit";
import {
  MessageTypeTunnelLog,
  MessageTypeTunnelStop,
} from "../../constants/websocket";
import Message from "../../interfaces/Message";
import Log from "./interfaces/Log";

interface LogEntries {
  [key: string]: Log[];
}

interface LogsState {
  tunnelLogs: LogEntries;
  communicationLogs: Log[];
}

const initialState: LogsState = { tunnelLogs: {}, communicationLogs: [] };

export const logsReducer = createReducer(initialState, {
  "REDUX_WEBSOCKET::MESSAGE": (
    state: LogsState,
    action: {
      meta: {
        timestamp: Date;
      };
      payload: { message: any };
    }
  ) => {
    // Tunnel logs
    if (action.payload.message.type === MessageTypeTunnelLog) {
      if (state.tunnelLogs[action.payload.message.tunnelId])
        state.tunnelLogs[action.payload.message.tunnelId].push({
          ...action.payload.message,
          timestamp: action.meta.timestamp,
        });
      else
        state.tunnelLogs[action.payload.message.tunnelId] = [
          {
            ...action.payload.message,
            timestamp: action.meta.timestamp,
          },
        ];
    }
    if (action.payload.message.type === MessageTypeTunnelStop)
      delete state.tunnelLogs[action.payload.message.tunnelId];

    state.communicationLogs.push({
      message: JSON.stringify(action.payload.message),
      class: action.payload.message.class,
      timestamp: action.meta.timestamp,
    });
  },
  "REDUX_WEBSOCKET::SEND": (
    state: LogsState,
    action: {
      meta: {
        timestamp: Date;
      };
      payload: Message<any>;
    }
  ) => {
    // Communication logs
    state.communicationLogs.push({
      message: JSON.stringify(action.payload),
      class: "info",
      timestamp: action.meta.timestamp,
    });
  },
});

export default logsReducer;
