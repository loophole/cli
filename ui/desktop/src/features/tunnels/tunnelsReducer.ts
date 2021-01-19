import {
  DEFAULT_PREFIX,
  WEBSOCKET_MESSAGE,
  WEBSOCKET_SEND,
} from "@giantmachines/redux-websocket";
import { createReducer } from "@reduxjs/toolkit";
import {
  MessageTypeLoadingFailure,
  MessageTypeLoadingStart,
  MessageTypeLoadingSuccess,
  MessageTypeTunnelStart,
  MessageTypeTunnelStartFailure,
  MessageTypeTunnelStartSuccess,
  MessageTypeTunnelStop,
  PrefixMessageTypeTunnelStart,
} from "../../constants/websocket";
import ExposeDirectoryMessage from "../../interfaces/ExposeDirectoryMessage";
import ExposeHttpPortMessage from "../../interfaces/ExposeHttpPortMessage";
import Message from "../../interfaces/Message";
import { MessageTypeDeleteFailedTunnel } from "./actions";
import Tunnel from "./interfaces/Tunnel";
import TunnelsState from "./interfaces/TunnelsState";

const initialState: TunnelsState = { tunnels: [] };

export const tunnelsReducer = createReducer(initialState, {
  [`${DEFAULT_PREFIX}::${WEBSOCKET_MESSAGE}`]: (
    state: TunnelsState,
    action
  ) => {
    switch (action.payload.message.type) {
      case MessageTypeTunnelStart: {
        if (action.payload.message.tunnelId) {
          const tunnelIndex = state.tunnels.findIndex(
            (tunnel: Tunnel) =>
              tunnel.tunnelId === action.payload.message.tunnelId
          );
          state.tunnels[tunnelIndex] = {
            ...state.tunnels[tunnelIndex],
            ...action.payload.message,
            type: state.tunnels[tunnelIndex].type,
          };
        } else {
          state.tunnels.push(action.payload.message);
        }
        break;
      }
      case MessageTypeTunnelStartSuccess: {
        if (action.payload.message.tunnelId) {
          const tunnelIndex = state.tunnels.findIndex(
            (tunnel: Tunnel) => tunnel.tunnelId === action.payload.message.tunnelId
          );
          state.tunnels[tunnelIndex] = {
            ...state.tunnels[tunnelIndex],
            ...action.payload.message,
            startTime: action.meta.timestamp,
            started: true,
            loading: false,
            loadingMsg: "",
            type: state.tunnels[tunnelIndex].type,
          };
        } else {
          state.tunnels.push(action.payload.message);
        }
        break;
      }
      case MessageTypeTunnelStartFailure: {
        const tunnelIndex = state.tunnels.findIndex(
          (tunnel: Tunnel) =>
            tunnel.tunnelId === action.payload.message.tunnelId
        );
        if (tunnelIndex === -1) break;
        state.tunnels[tunnelIndex] = {
          ...state.tunnels[tunnelIndex],
          loading: false,
          loadingMsg: "",
          error: true,
          errorMsg: action.payload.message.error,
        };
        break;
      }
      case MessageTypeTunnelStop: {
        state.tunnels = state.tunnels.filter(
          (tunnel) => tunnel.tunnelId !== action.payload.message.tunnelId
        );
        break;
      }
      case MessageTypeLoadingStart: {
        let tunnelIndex = state.tunnels.findIndex(
          (tunnel: Tunnel) => tunnel.tunnelId === action.payload.message.tunnelId
        );
        if (tunnelIndex === -1) break;

        state.tunnels[tunnelIndex] = {
          ...state.tunnels[tunnelIndex],
          loading: true,
          loadingMsg: action.payload.message.message,
        };
        break;
      }
      case MessageTypeLoadingSuccess: {
        let tunnelIndex = state.tunnels.findIndex(
          (tunnel: Tunnel) => tunnel.tunnelId === action.payload.message.tunnelId
        );
        if (tunnelIndex === -1) break;

        state.tunnels[tunnelIndex] = {
          ...state.tunnels[tunnelIndex],
          loading: false,
          loadingMsg: "",
        };
        break;
      }
      case MessageTypeLoadingFailure: {
        let tunnelIndex = state.tunnels.findIndex(
          (tunnel: Tunnel) => tunnel.tunnelId === action.payload.message.tunnelId
        );
        if (tunnelIndex === -1) break;

        state.tunnels[tunnelIndex] = {
          ...state.tunnels[tunnelIndex],
          loading: false,
          error: true,
          loadingMsg: "",
          errorMsg: action.payload.message.error,
        };
        break;
      }
    }
  },
  [`${DEFAULT_PREFIX}::${WEBSOCKET_SEND}`]: (
    state: TunnelsState,
    action: {
      payload: Message<ExposeHttpPortMessage> | Message<ExposeDirectoryMessage>;
    }
  ) => {
    if (action.payload.type.startsWith(PrefixMessageTypeTunnelStart)) {
      state.tunnels.push({
        tunnelId: action.payload.payload.remote.tunnelId,
        siteId: action.payload.payload.remote.siteId,
        type: action.payload.type.replace(PrefixMessageTypeTunnelStart, ""),
        started: false,
        usingBasicAuth: !!action.payload.payload.remote.basicAuthUsername,
        basicAuthUsername: action.payload.payload.remote.basicAuthUsername,
        basicAuthPassword: action.payload.payload.remote.basicAuthPassword,
        proxyErrorDisabled: action.payload.payload.remote.disableProxyErrorPage,
      });
    }
  },
  [MessageTypeDeleteFailedTunnel]: (
    state: TunnelsState,
    action: {
      payload: { tunnel: Tunnel };
    }
  ) => {
    console.log(action, action.payload, action.payload.tunnel);
    state.tunnels = state.tunnels
      .filter((tunnel) => tunnel.tunnelId !== action.payload.tunnel.tunnelId);
  },
});

export default tunnelsReducer;
