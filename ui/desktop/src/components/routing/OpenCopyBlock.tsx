import React, { useState } from "react";
import { useDispatch } from "react-redux";
import { send } from "@giantmachines/redux-websocket";
import Message from "../../interfaces/Message";
import { CopyToClipboard } from "react-copy-to-clipboard";
import OpenInBrowserMessage from "../../interfaces/OpenInBrowserMessage";
import { MessageTypeOpenInBrowser } from "../../constants/websocket";

const OpenCopyBlock = (props: any) => {
  const dispatch = useDispatch();
  const [copied, setCopied] = useState(false);

  let { target, targetName, open, copy } = props;
  if (typeof open === 'undefined') open = true;
  if (typeof copy === 'undefined') copy = true;

  if (!targetName) targetName = target;

  const navigate = () => {
    const message: Message<OpenInBrowserMessage> = {
      type: MessageTypeOpenInBrowser,
      payload: {
        url: target,
      },
    };

    dispatch(send(message));
  };

  const onCopiedEffect = () => {
    setCopied(true);
    setTimeout(() => {
      setCopied(false);
    }, 1000);
  };

  return (
    <pre className="mt-4 mb-4">
      <code>
        {targetName}{" "}
        {open ? (
          <span
            className="tag is-primary is-clickable"
            onClick={navigate}
            onKeyDown={navigate}
          >
            Open
          </span>
        ) : null}
        {open && copy ? " " : ""}
        {copy ? (
          <CopyToClipboard text={target} onCopy={onCopiedEffect}>
            <span className="tag is-info is-clickable">
              {copied ? "Copied!" : "Copy"}
            </span>
          </CopyToClipboard>
        ) : null}
      </code>
    </pre>
  );
};

export default OpenCopyBlock;
