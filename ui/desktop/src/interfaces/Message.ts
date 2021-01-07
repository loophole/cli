import OpenInBrowserMessage from './OpenInBrowserMessage';
import StartTunnelMessage from './StartTunnelMessage';
import StopTunnelMessage from './StopTunnelMessage';

export default interface Message {
    messageType: string;
    startTunnelMessage?: StartTunnelMessage;
    stopTunnelMessage?: StopTunnelMessage;
    openInBrowserMessage?: OpenInBrowserMessage
}