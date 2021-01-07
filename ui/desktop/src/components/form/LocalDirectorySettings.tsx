import React from "react";
import classNames from 'classnames';
interface LocalDirectorySettingsProps {
  pathValue: string;
  pathChangeCallback: Function;
}
const LocalDirectorySettings = (props: LocalDirectorySettingsProps): JSX.Element => {
  const path = props.pathValue;
  const setPath = props.pathChangeCallback;

  const isPathValid = (): boolean => {
    return path.length >= 1;
  }

  return (
    <div>
      <div className="field">
        <label className="label">Path</label>
        <div className="control has-icons-left has-icons-right">
          <input
            className={classNames({
              input: true,
              "is-success": isPathValid(),
              "is-danger": !isPathValid(),
            })}
            type="text"
            placeholder="Host on which the server is running"
            value={path}
            onChange={(e) => setPath(e.target.value)}
          />
          <span className="icon is-small is-left">
            <i className="fas fa-signature"></i>
          </span>
          <span className="icon is-small is-right">
            <i 
                    className={classNames({
                      fas: true,
                      "fa-check": isPathValid(),
                      "fa-exclamation-triangle": !isPathValid(),
                    })}></i>
          </span>
        </div>
              {isPathValid() ? (
                <p className="help is-success">Path is valid</p>
              ) : (
                <p className="help is-danger">Path is invalid</p>
              )}
      </div>
    </div>
  );
};

export default LocalDirectorySettings;
