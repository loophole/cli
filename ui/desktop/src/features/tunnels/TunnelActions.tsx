import React, { useState } from "react";
import { useHistory } from "react-router-dom";
import StopTunnelModal from "./StopTunnelModal";

const TunnelActions = (props: any) => {
  const history = useHistory();
  const [modalVisible, setModalVisible] = useState(false);

  return (
    <div className="card">
    <StopTunnelModal hideAction={() => { setModalVisible(false); history.push("/tunnels")}} visible={modalVisible} tunnel={props.tunnel} />
      <header className="card-header">
        <p className="card-header-title">Actions</p>
      </header>
      <div className="card-content has-text-right">
        <div className="content">
          <button
            className="button is-danger"
            onClick={() => { setModalVisible(true); }}
            onKeyDown={() => { setModalVisible(true); }}
          >
            <span className="icon">
              <i className="fa fa-times-circle" />
            </span>
            <span>{props.tunnel.error ? "Delete" : "Stop"} the tunnel</span>
          </button>
        </div>
      </div>
    </div>
  );
};

export default TunnelActions;
