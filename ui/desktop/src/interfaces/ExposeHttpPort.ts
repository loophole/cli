import DisplayOptions from './DisplayOptions';
import LocalHTTPEndpointSpecs from './LocalHTTPEndpointSpecs';
import RemoteEndpointSpecs from './RemoteEndpointSpecs';

export default interface ExposeHttpPort {
	local:   LocalHTTPEndpointSpecs;
	remote:  RemoteEndpointSpecs;
	display: DisplayOptions;
}