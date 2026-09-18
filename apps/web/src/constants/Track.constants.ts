export const TRACK_RANGE_PRESETS = {
	0: {
		course: [20, 20, 20],
		crosstrack: [255, 255, 255],
	},
	1: {
		course: [20, 20, 20],
		crosstrack: [10, 10, 10],
	},
	2: {
		course: [20, 10, 5],
		crosstrack: [6, 6, 6],
	},
	3: {
		course: [5, 10, 20],
		crosstrack: [5, 3, 3],
	},
	4: {
		course: [5, 5, 5],
		crosstrack: [3, 3, 3],
	},
};

export const STATIC_MODE = {
	OFF: 0,
	AIN_NORMAL: 1,
	AIN_INVERSE: 2,
	PROGRAM: 3,
	ACCEL: 4,
} as const;

export const TRACK_STATIC_MODE_BASE = [
	{
		id: STATIC_MODE.OFF,
		value: 'label.turned-off',
	},
	{
		id: STATIC_MODE.AIN_NORMAL,
		value: 'label.by-ignition',
	},
	{
		id: STATIC_MODE.ACCEL,
		value: 'label.by-accelerometer',
	},
];

export const TRACK_STATIC_MODE_2 = [
	{
		id: STATIC_MODE.OFF,
		value: 'label.turned-off',
	},
	{
		id: STATIC_MODE.ACCEL,
		value: 'label.by-accelerometer',
	},
];

export const TRACK_ENTRY_STATUS_STATIC_MODE = [
	{
		id: STATIC_MODE.AIN_NORMAL,
		value: 'message.logical-zero',
	},
	{
		id: STATIC_MODE.AIN_INVERSE,
		value: 'message.logical-one',
	},
];

export const TRACK_STATIC_MODE_BY_SPEED = [
	{
		id: STATIC_MODE.OFF,
		value: 'label.turned-off',
	},
	{
		id: STATIC_MODE.AIN_NORMAL,
		value: 'label.by-ignition',
	},
	{
		id: STATIC_MODE.PROGRAM,
		value: 'label.by-speed',
	},
];

export const TRACK_ALTERNATIVE_MODE = [
	{
		id: 0,
		value: 'message.track-alternative-mode-gnss-only',
	},
	{
		id: 1,
		value: 'message.track-alternative-mode-static-substitution',
	},
	{
		id: 2,
		value: 'message.track-alternative-mode-gnss-priority',
	},
	{
		id: 3,
		value: 'message.track-alternative-mode-light-alt-priority',
	},
	{
		id: 4,
		value: 'message.track-alternative-mode-full-alt-priority',
	},
	{
		id: 5,
		value: 'message.track-alternative-mode-soft-validation-with-substitution',
	},
	{
		id: 6,
		value: 'message.track-alternative-mode-hard-validation-with-substitution',
	},
	{
		id: 7,
		value: 'message.track-alternative-mode-soft-validation-without-substitution',
	},
	{
		id: 8,
		value: 'message.track-alternative-mode-hard-validation-without-substitution',
	},
];

export const TRACK_OPERATING_MODE_MINIMAL = [
	{
		id: 0,
		value: 'message.track-operating-mode-disabled',
	},
	{
		id: 1,
		value: 'message.track-operating-mode-yandex-locator',
	},
];

export const TRACK_OPERATING_MODE_ADVENCED = [
	{
		id: 0,
		value: 'message.track-operating-mode-disabled',
	},
	{
		id: 1,
		value: 'message.track-operating-mode-yandex-locator',
	},
	{
		id: 2,
		value: 'message.track-operating-mode-google-locator',
	},
];

export const TRACK_NAVTIMESYNC = [
	{
		id: 0,
		value: 'message.track-navigation-time-priority',
	},
	{
		id: 1,
		value: 'message.track-ntp-time-priority',
	},
	{
		id: 2,
		value: 'message.track-time-is-ntp-only',
	},
];
