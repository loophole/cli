import DisplayOptions from './DisplayOptions';
import LocalHTTPEndpointSpecs from './LocalHTTPEndpointSpecs';
import RemoteEndpointSpecs from './RemoteEndpointSpecs';

export default interface ExposeHttpPortMessage {
	local:   LocalHTTPEndpointSpecs;
	remote:  RemoteEndpointSpecs;
}