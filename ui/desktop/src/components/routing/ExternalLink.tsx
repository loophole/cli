import React from "react";
import { useDispatch } from "react-redux";
import { send } from "@giantmachines/redux-websocket";
import Message from "../../interfaces/Message";

const ExternalLink = (props: any) => {
  const dispatch = useDispatch();

  let { url, urlName } = props;
  if (!urlName) urlName = url;

  const navigate = () => {
    const message: Message = {
      messageType: "MT_OpenInBrowser",
      openInBrowserMessage: {
        url,
      },
    };

    dispatch(send(message));
  };
  return (
    <button className="button is-ghost" onClick={navigate} onKeyDown={navigate}>
      {urlName}
    </button>
  );
};

export default ExternalLink;