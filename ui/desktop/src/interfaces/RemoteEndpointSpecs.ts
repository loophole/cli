import Endpoint from "./Endpoint";

export default interface RemoteEndpointSpecs {
  gatewayEndpoint?: Endpoint;
  apiEndpoint?: Endpoint;
  identityFile?: string;
  siteId?: string;
  tunnelId: string;
  basicAuthUsername?: string;
  basicAuthPassword?: string;
  disableProxyErrorPage: boolean;
}
