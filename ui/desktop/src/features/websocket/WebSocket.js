import { useEffect } from "react";
import { useDispatch } from "react-redux";
import { connect, disconnect } from '@giantmachines/redux-websocket';

export const WebSocket = () => {
  const dispatch = useDispatch();

  const host = `${window.location.protocol === "http:" ? "ws" : "wss"}://${
    window.location.host
  }/ws`;

  useEffect(() => {
    dispatch(connect(host));

    return () => {
      dispatch(disconnect());
    };
  });

  return null;
};
