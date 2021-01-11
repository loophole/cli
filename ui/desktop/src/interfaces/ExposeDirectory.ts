import DisplayOptions from './DisplayOptions';
import LocalDirectorySpecs from './LocalDirectorySpecs';
import RemoteEndpointSpecs from './RemoteEndpointSpecs';

export default interface ExposeDirectory {
	local:   LocalDirectorySpecs;
	remote:  RemoteEndpointSpecs;
	display: DisplayOptions;
}