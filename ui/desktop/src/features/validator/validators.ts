const validIpAddressRegex =
  "^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]).){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$";
const validHostnameRegex =
  "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]).)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9-]*[A-Za-z0-9])$";

export const isLoopholeHostnameValid = (hostname: string): boolean => {
  return (
    hostname.match(/^[a-z](?:-?[a-z0-9])*$/) !== null &&
    hostname.length > 1 &&
    hostname.length < 30
  );
};

export const isLocalPathValid = (path: string): boolean => {
  return path.length > 0;
};

export const isLocalHostValid = (address: string): boolean => {
  return (
    address.match(validIpAddressRegex) !== null ||
    address.match(validHostnameRegex) !== null
  );
};

export const isLocalPortValid = (port: number): boolean => {
  return port > 0 && port <= 65535;
};

export const isBasicAuthUsernameValid = (
  basicAuthUsername: string
): boolean => {
  return basicAuthUsername.length >= 3;
};
export const isBasicAuthPasswordValid = (
  basicAuthPassword: string
): boolean => {
  return basicAuthPassword.length >= 3;
};
