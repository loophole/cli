import React from "react";
import { useSelector } from "react-redux";

const AboutApplication = () => {
  const appState = useSelector((store: any) => store.config);

  return (
    <div className="content has-text-right">
      <p>
        <small className="is-uppercase">
          Version: {appState.version} ({appState.commitHash})
        </small>
      </p>
    </div>
  );
};

export default AboutApplication;
