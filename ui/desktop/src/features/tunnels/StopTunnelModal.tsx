import React from "react";
import { send } from "@giantmachines/redux-websocket";
import classNames from "classnames";
import Message from "../../interfaces/Message";
import Tunnel from "./interfaces/Tunnel";
import { MessageTypeRequestTunnelStop } from "../../constants/websocket";
import StopTunnelMessage from "../../interfaces/StopTunnelMessage";
import { useDispatch } from "react-redux";
import { MessageTypeDeleteFailedTunnel } from "./actions";

interface StopTunnelProps {
  tunnel: Tunnel;
  visible: boolean;
  hideAction: Function;
}

const StopTunnelModal = (props: StopTunnelProps) => {
  const dispatch = useDispatch();
  const stopTunnel = () => {
    if (!props.tunnel.siteId) return;

    const message: Message<StopTunnelMessage> = {
      type: MessageTypeRequestTunnelStop,
      payload: {
        tunnelId: props.tunnel.tunnelId,
      },
    };

    dispatch(send(message));

    props.hideAction();
  };
  const deleteFailedTunnel = () => {
    if (!props.tunnel.error) return;

    const action = {
      type: MessageTypeDeleteFailedTunnel,
      payload: {
        tunnel: props.tunnel,
      },
    };
    console.log(action);
    dispatch(action);
    props.hideAction();
  };
  return (
    <div
      className={classNames({
        modal: true,
        "is-active": props.visible,
        "is-clipped": props.visible,
      })}
    >
      <div className="modal-background"></div>
      <div className="modal-card">
        <header className="modal-card-head">
          <p className="modal-card-title">Shut down the tunnel</p>
          <button
            className="delete"
            aria-label="close"
            onClick={() => props.hideAction()}
            onKeyDown={() => props.hideAction()}
          ></button>
        </header>
        <section className="modal-card-body">
          Are you sure you want to stop tunnel {props.tunnel.siteId}?
        </section>
        <footer className="modal-card-foot">
          <button
            className="button is-danger"
            onClick={!props.tunnel.error ? stopTunnel : deleteFailedTunnel}
            onKeyDown={!props.tunnel.error ? stopTunnel : deleteFailedTunnel}
          >
            Yes, {props.tunnel.error ? "delete" : "stop"} it
          </button>
          <button
            className="button"
            onClick={() => props.hideAction()}
            onKeyDown={() => props.hideAction()}
          >
            No, take me back
          </button>
        </footer>
      </div>
    </div>
  );
};

export default StopTunnelModal;
