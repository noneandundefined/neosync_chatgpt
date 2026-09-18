export const STEP_KEYS = ['device', 'sim', 'server', 'events', 'track', 'inputs', 'outputs', 'rs485', 'bluetooth', 'onewire'] as const;

export type StepKey = (typeof STEP_KEYS)[number];
