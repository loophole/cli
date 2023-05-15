import React from "react";

interface DirectorySettingsProps {
  usingValue: boolean;
  usingChangeCallback: Function;
}

const DirectorySettings = (props: DirectorySettingsProps): JSX.Element => {
  const disableDirectoryListing = props.usingValue;
  const setDisableDirectoryListing = props.usingChangeCallback;

  return (
    <div>
      <div className="field">
        <div className="control">
          <label className="checkbox">
            <input
              type="checkbox"
              onChange={(e) => {
                setDisableDirectoryListing(!disableDirectoryListing);
              }}
            />{" "}
            I want to disable Directory Listing
          </label>
        </div>
      </div>
    </div>
  );
};

export default DirectorySettings;
