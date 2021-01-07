import Endpoint from "./Endpoint";

export default interface RemoteEndpointSpecs {
  gatewayEndpoint?: Endpoint;
  apiEndpoint?: Endpoint;
  identityFile?: string;
  siteID?: string;
  basicAuthUsername?: string;
  basicAuthPassword?: string;
}
