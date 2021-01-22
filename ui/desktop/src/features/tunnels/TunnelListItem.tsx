import classNames from "classnames";
import React, { useState } from "react";
import { useHistory } from "react-router-dom";
import Tunnel from "./interfaces/Tunnel";
import StopTunnelModal from "./StopTunnelModal";
import TunnelStatusIcons from "./TunnelStatusIcon";

interface TunnelProps {
  tunnel: Tunnel;
}

const TunnelListItem = (props: TunnelProps) => {
  const history = useHistory();
  const [modalVisible, setModalVisible] = useState(false);

  const navigateToInfo = () => {
    history.push(`/tunnels/${props.tunnel.tunnelId}`);
  };

  const showConfirmationModal = (
    e:
      | React.MouseEvent<HTMLButtonElement>
      | React.KeyboardEvent<HTMLButtonElement>
  ) => {
    e.stopPropagation();

    setModalVisible(true);
  };

  return (
    <div className="column is-half is-narrow">
      <StopTunnelModal
        hideAction={() => setModalVisible(false)}
        visible={modalVisible}
        tunnel={props.tunnel}
      />
      <div
        className={classNames({
          notification: true,
          "is-success": props.tunnel.started,
          "is-info": !props.tunnel.started && !props.tunnel.error,
          "is-danger": props.tunnel.error,
          "is-clickable": true,
        })}
        style={{
          height: "200px",
        }}
        onClick={navigateToInfo}
        onKeyDown={navigateToInfo}
      >
        {props.tunnel.siteId ? (
          <button
            className="delete"
            onClick={showConfirmationModal}
            onKeyDown={showConfirmationModal}
            title="Delete tunnel"
          ></button>
        ) : null}
        <div className="content">
          <p>
            <strong
              style={{
                overflowWrap: "break-word",
              }}
            >
              {props.tunnel.loading ? (
                <span>{props.tunnel.loadingMsg}</span>
              ) : null}
              {props.tunnel.error ? <span>{props.tunnel.errorMsg}</span> : null}
              {props.tunnel.started ? (
                <span>
                  {props.tunnel.siteId}
                  <small>
                    .loophole.site{" "}
                    {props.tunnel.localAddr ? (
                      <span>&rarr; {props.tunnel.localAddr}</span>
                    ) : null}
                  </small>
                </span>
              ) : null}
            </strong>
          </p>

          <p>
            {props.tunnel.started ? (
              <small>
                <span className="icon">
                  <i className="fas fa-check" />
                </span>
                {props.tunnel.startTime
                  ? `Running since ${props.tunnel.startTime.toLocaleString()}`
                  : "Running"}
              </small>
            ) : null}
            {props.tunnel.error ? (
              <small>
                <span className="icon">
                  <i className="fas fa-times" />
                </span>
                Starting tunnel failed. Please try again
              </small>
            ) : null}
            {!props.tunnel.started && !props.tunnel.error ? (
              <small>
                <span className="icon">
                  <i className="fas fa-circle-notch fa-spin" />
                </span>
                Starting up...
              </small>
            ) : null}
          </p>
        </div>
        {props.tunnel.started ? (
          <TunnelStatusIcons tunnel={props.tunnel} />
        ) : null}
      </div>
    </div>
  );
};

export default TunnelListItem;
