import React from "react";

import { Beforeunload } from "react-beforeunload";
import { useSelector } from "react-redux";

const EventCatcher = (props) => {
  const tunnelsState = useSelector((store) => store.tunnels);

  return (
    <Beforeunload
      onBeforeunload={(event) =>
        tunnelsState.tunnels.length ? event.preventDefault() : null
      }
    >
      {props.children}
    </Beforeunload>
  );
};

export default EventCatcher;
