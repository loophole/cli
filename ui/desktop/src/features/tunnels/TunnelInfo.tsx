import React from "react";
import QRCode from "qrcode.react";

import OpenCopyBlock from "../../components/routing/OpenCopyBlock";
import Tunnel from "./interfaces/Tunnel";

interface TunnelInfoProps {
  tunnel: Tunnel;
}

const TunnelInfo = (props: TunnelInfoProps) => {
  const { tunnel } = props;

  const downloadQR = () => {
    const canvas: any = document.querySelector(".qrcode-container > canvas");
    if (!canvas) return;
    const pngUrl = canvas
      .toDataURL("image/png")
      .replace("image/png", "image/octet-stream");
    let downloadLink = document.createElement("a");
    downloadLink.href = pngUrl;
    downloadLink.download = `${tunnel.siteId}-qrcode.png`;
    document.body.appendChild(downloadLink);
    downloadLink.click();
    document.body.removeChild(downloadLink);
  };

  return (
    <div className="content">
      <div className="has-text-centered">
        {tunnel.started && tunnel.siteAddrs && tunnel.siteAddrs[0] ? (
          <div className="qrcode-container">
            <QRCode
              value={tunnel.siteAddrs[0]}
              renderAs="canvas"
              width="80%"
              size={256}
              level="H"
            />
            <div>
              <button
                className="button is-ghost"
                onClick={downloadQR}
                onKeyDown={downloadQR}
              >
                Download QRCode
              </button>
            </div>
          </div>
        ) : null}
        {!tunnel.started && tunnel.loading ? tunnel.loadingMsg : null}
        {!tunnel.started && tunnel.error ? tunnel.errorMsg : null}
      </div>
      {tunnel.started ? (
        tunnel.siteAddrs && tunnel.siteAddrs.length > 1 ? (
          <ul>
            {tunnel.siteAddrs.map((addr) => (
              <li>
                <OpenCopyBlock target={addr} />
              </li>
            ))}
          </ul>
        ) : (
          <span>
            {tunnel.siteAddrs && tunnel.siteAddrs[0] ? (
              <OpenCopyBlock target={tunnel.siteAddrs[0]} />
            ) : null}
          </span>
        )
      ) : null}
    </div>
  );
};

export default TunnelInfo;
