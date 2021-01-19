import { useSelector } from "react-redux";

import classNames from "classnames";

const ProfilePage = () => {
  const appState = useSelector((state: any) => state.config);

  return (
    <div className="container">
      <h1 className="subtitle is-4">Profile</h1>
      <hr />
      <div className="context-box">
        <article className="media">
          <figure className="media-left">
            <p className="image is-64x64">
              <img src={appState.user.picture} alt="User avatar" />
            </p>
          </figure>
          <div className="media-content">
            <div className="content">
              <p>
                <strong>{appState.user.name}</strong>{" "}
                <small>@{appState.user.nickname}</small>{" "}
              </p>
              <p>
                <div className="tags has-addons">
                  <span className="tag is-black">Email</span>
                  <span className="tag">{appState.user.email}</span>
                  <span
                    className={classNames({
                      tag: true,
                      "is-light": !appState.user.email_verified,
                      "is-success": appState.user.email_verified,
                    })}
                  >
                    {appState.user.email_verified ? "Verified" : "Unverified"}
                  </span>
                </div>
              </p>
            </div>
          </div>
        </article>
      </div>
    </div>
  );
};

export default ProfilePage;
