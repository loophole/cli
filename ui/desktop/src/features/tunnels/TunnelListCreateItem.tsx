import React from "react";
import { useHistory } from "react-router";

interface TunnelCreateProps {}

const TunnelListItem = (props: TunnelCreateProps) => {
  const history = useHistory();
  const navigate = () => {
    history.push("/tunnels/create");
  };

  return (
    <div className="column is-half is-narrow">
      <div
        className="notification is-light is-clickable is-half is-flex is-align-items-center is-justify-content-center"
        style={{
          height: "200px",
        }}
        onClick={navigate}
      >
        <span className="icon is-large">
          <i className="fas fa-plus fa-4x" />
        </span>
      </div>
    </div>
  );
};

export default TunnelListItem;
