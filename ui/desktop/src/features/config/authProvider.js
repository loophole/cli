import { useContext, createContext, useState } from "react";

const authContext = createContext();

const AuthProvider = ({ children }) => {
  const auth = useProvideAuth();
  return <authContext.Provider value={auth}>{children}</authContext.Provider>;
};

export const useAuth = () => {
  return useContext(authContext);
};

const useProvideAuth = () => {
  const [loggedIn, setLoggedIn] = useState(null);

  const login = (cb) => {
    setLoggedIn(true);
    cb();
  };
  const logout = (cb) => {
    setLoggedIn(false);
    cb();
  };

  return {
    loggedIn,
    login,
    logout
  };
};

export default AuthProvider;