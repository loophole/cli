import { MessageType } from "../constants/websocket";

export default interface Message<T> {
    type: MessageType;
    payload: T;
}
