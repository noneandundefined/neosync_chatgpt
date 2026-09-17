import ConfigSim from './ConfigurationBlocks/Sim';
import ConfigTrack from './ConfigurationBlocks/Track';
import ConfigRS485 from './ConfigurationBlocks/RS485';
import ConfigWire from './ConfigurationBlocks/OneWire';
import ConfigInput from './ConfigurationBlocks/Inputs';
import ConfigDevice from './ConfigurationBlocks/Device';
import ConfigEvents from './ConfigurationBlocks/Events';
import ConfigOutput from './ConfigurationBlocks/Output';
import ConfigServer from './ConfigurationBlocks/Server';
import ConfigBluetooth from './ConfigurationBlocks/Bluetooth';

interface ConfigurationStepperProps {
	currentStep: number;
}

const ConfigurationStepper: React.FC<ConfigurationStepperProps> = ({ currentStep }) => {
	switch (currentStep) {
		case 0:
			return <ConfigDevice />;
		case 1:
			return <ConfigSim />;
		case 2:
			return <ConfigServer />;
		case 3:
			return <ConfigEvents />;
		case 4:
			return <ConfigTrack />;
		case 5:
			return <ConfigInput />;
		case 6:
			return <ConfigOutput />;
		case 7:
			return <ConfigRS485 />;
		case 8:
			return <ConfigBluetooth />;
		case 9:
			return <ConfigWire />;
	}
};

export default ConfigurationStepper;
