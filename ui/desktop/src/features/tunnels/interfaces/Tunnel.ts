export default interface Tunnel {
  siteId?: string;
  tunnelId: string;
  type: string;
  startTime?: Date;
  started: boolean;
  usingBasicAuth: boolean;
  basicAuthUsername?: string;
  basicAuthPassword?: string;
  proxyErrorDisabled: boolean;
  localAddr?: string;
  siteAddrs?: string[];

  loading?: boolean;
  loadingMsg?: string;
  error?: boolean;
  errorMsg?: string;
}
