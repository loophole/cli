import React from "react";
import { useSelector } from "react-redux";
import Tunnel from "../features/tunnels/interfaces/Tunnel";
import TunnelListCreateItem from "../features/tunnels/TunnelListCreateItem";
import TunnelListItem from "../features/tunnels/TunnelListItem";

const TunnelsPage = (): JSX.Element => {
  const tunnelsState = useSelector((store: any) => store.tunnels);

  const tunnelBoxes = tunnelsState.tunnels.map((tunnel: Tunnel) => {
    return (
      <TunnelListItem
        key={tunnel.siteId}
        tunnel={tunnel}
      />
    );
  });

  return (
    <div className="container">
      <h1 className="subtitle is-4">
        Select a tunnel to manage or create a new one.
      </h1>
      <hr />
      <div className="context-box">
        <div className="columns is-multiline">
          {tunnelBoxes}
          <TunnelListCreateItem />
        </div>
      </div>
    </div>
  );
};

export default TunnelsPage;
