import React from "react";

import { useSelector } from "react-redux";
import OpenCopyBlock from "../components/routing/OpenCopyBlock";

const FeedbackPage = (): JSX.Element => {
  const appState = useSelector((store: any) => store.config);

  return (
    <div className="container">
      <h1 className="subtitle is-4">Feedback</h1>
      <hr />
      <div className="context-box">
        {appState.feedbackFormUrl ? (
          <div>
            Our feedback form is available as Google Form and can be accessed by
            using the link:
            <OpenCopyBlock target={appState.feedbackFormUrl} />
            Please share your thoughs on <em>Loophole</em> there.
          </div>
        ) : (
          <div>
            Feedback options is not avable at the moment, please bear with us
            and check out next time.
          </div>
        )}
      </div>
    </div>
  );
};

export default FeedbackPage;
