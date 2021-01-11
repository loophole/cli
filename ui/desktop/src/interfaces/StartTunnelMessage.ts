import ExposeHttpPort from "./ExposeHttpPort";
import ExposeDirectory from "./ExposeDirectory";

export default interface StartTunnelMessage {
    tunnelType: string;
    exposeHttpConfig?: ExposeHttpPort;
    exposeDirectoryConfig?: ExposeDirectory;
    exposeWebdavConfig?: ExposeDirectory;
}