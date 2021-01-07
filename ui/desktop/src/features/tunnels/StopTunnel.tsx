import { send } from "@giantmachines/redux-websocket";
import React from "react";
import { useDispatch } from "react-redux";
import  Message from "../../interfaces/Message"

const StopTunnel = (props: any) => {
  const dispatch = useDispatch();
  const { tunnel } = props;

  const stopTunnel = () => {
    if (!tunnel) return;

    const message : Message = {
      messageType: "MT_StopTunnel",
      stopTunnelMessage: {
        siteId: tunnel.siteId
      }
    }

    dispatch(send(message));
  };

  return (
    <div className="card ">
      <header className="card-header">
        <p className="card-header-title has-background-danger-light has-text-danger">Danger Zone</p>
      </header>
      <div className="card-content has-text-centered">
        <div className="content">
          <button
            className="button is-danger"
            onClick={stopTunnel}
            onKeyDown={stopTunnel}
          >
            Stop the tunnel
          </button>
        </div>
      </div>
    </div>
  );
};

export default StopTunnel;
