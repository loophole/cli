export type MessageType = string;

export const MessageTypeLog: MessageType = "MT_Log";

export const MessageTypeTunnelLog: MessageType = "MT_TunnelLog";

export const MessageTypeAppStart: MessageType = "MT_ApplicationStart";
export const MessageTypeAppStop: MessageType = "MT_ApplicationStop";
export const MessageTypeNewVersionAvailable: MessageType =
  "MT_ApplicationNewVersionAvailable";

export const MessageTypeLogin: MessageType = "MT_Login";
export const MessageTypeLoginSuccess: MessageType = "MT_LoginSuccess";
export const MessageTypeLoginFailure: MessageType = "MT_LoginFailure";

export const MessageTypeLogout: MessageType = "MT_Logout";
export const MessageTypeLogoutSuccess: MessageType = "MT_LogoutSuccess";
export const MessageTypeLogoutFailure: MessageType = "MT_LogoutFailure";

export const MessageTypeTunnelStart: MessageType = "MT_TunnelStart";
export const MessageTypeTunnelStartSuccess: MessageType =
  "MT_TunnelStartSuccess";
export const MessageTypeTunnelStartFailure: MessageType =
  "MT_TunnelStartFailure";

export const MessageTypeTunnelStop: MessageType = "MT_TunnelStop";

export const MessageTypeLoadingStart: MessageType = "MT_LoadingStart";
export const MessageTypeLoadingSuccess: MessageType = "MT_LoadingSuccess";
export const MessageTypeLoadingFailure: MessageType = "MT_LoadingFailure";

export const MessageTypeOpenInBrowser: MessageType = "MT_OpenInBrowser";

export const MessageTypeRequestTunnelStop: MessageType = "MT_RequestTunnelStop";

export const PrefixMessageTypeTunnelStart: string = "MT_RequestTunnelStart_";
export const MessageTypeRequestTunnelStartHTTP: MessageType = `${PrefixMessageTypeTunnelStart}HTTP`;
export const MessageTypeRequestTunnelStartDirectory: MessageType =  `${PrefixMessageTypeTunnelStart}Directory`;
export const MessageTypeRequestTunnelStartWebDav: MessageType = `${PrefixMessageTypeTunnelStart}WebDav`;

export const MessageTypeRequestLogout: MessageType = "MT_RequestLogout";
export const MessageTypeRequestLogin: MessageType = "MT_RequestLogin";