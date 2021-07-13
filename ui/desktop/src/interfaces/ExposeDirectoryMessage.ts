import DisplayOptions from './DisplayOptions';
import LocalDirectorySpecs from './LocalDirectorySpecs';
import RemoteEndpointSpecs from './RemoteEndpointSpecs';

export default interface ExposeDirectoryMessage {
	local:   LocalDirectorySpecs;
	remote:  RemoteEndpointSpecs;
	deactivatedirectorylisting: boolean;
}