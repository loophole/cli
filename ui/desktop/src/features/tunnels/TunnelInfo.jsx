import React from "react";
import QRCode from "qrcode.react";

import ExternalLink from '../../components/routing/ExternalLink';

const TunnelInfo = (props) => {
  const { tunnel } = props;

  return (
    <div className="card">
      <header className="card-header">
        <p className="card-header-title">Information</p>
      </header>
      <div className="card-content has-text-centered">
        <InfoPanel tunnel={tunnel} />
      </div>
    </div>
  );
};

const InfoPanel = (props) => {
  const { siteAddrs } = props.tunnel;

  return (
    <div>
      <div className="card-image">
        <QRCode
          renderAs="svg"
          size={512}
          width="80%"
          level="H"
          value={siteAddrs[0]}
        />
      </div>
      <div className="content">
        <span>Visit your webpage at </span>
        {siteAddrs.length > 1 ? (
          <ul>
            {siteAddrs.map((addr) => (
              <li>
              <ExternalLink url={addr} />
              </li>
            ))}
          </ul>
        ) : (
          <span>
            <ExternalLink url={siteAddrs[0]} />
          </span>
        )}
      </div>
    </div>
  );
};

export default TunnelInfo;
