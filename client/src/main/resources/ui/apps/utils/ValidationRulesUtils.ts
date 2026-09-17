export type ValidationResult = {
	valid: boolean;
	message?: string;
	argv?: string;
};

export interface ValidationRule {
	validate: (value: string) => ValidationResult;
}

export const numberRange = (min: number, max: number): ValidationRule => ({
	validate: (value: string) => {
		const num = Number(value);
		if (Number.isNaN(num)) return { valid: false, message: 'message.rule-value-numbers' };
		if (num < min)
			return {
				valid: false,
				message: `message.rule-value-numbers`,
				argv: `${min}...${max}`,
			};
		if (num > max)
			return {
				valid: false,
				message: `message.rule-value-numbers`,
				argv: `${min}...${max}`,
			};
		return { valid: true };
	},
});

export const ruleOnlyLetters: ValidationRule = {
	validate: (value: string) => {
		if (!/^[A-Za-zА-Яа-я]+$/.test(value)) {
			return { valid: false, message: 'message.rule-value-letters' };
		}

		return { valid: true };
	},
};

export const ruleOnlyDigits: ValidationRule = {
	validate: (value: string) => {
		if (!/^\+?\d+$/.test(value)) {
			return { valid: false, message: 'message.rule-value-numbers' };
		}

		return { valid: true };
	},
};

export const rulePassword: ValidationRule = {
	validate: (value: string) => {
		const regex = /^[A-Za-z0-9!"#$%&'()*+\-./:;<=>?@[\\\]^_{|}~]{0,8}$/;

		if (!regex.test(value)) {
			return { valid: false, message: 'message.requirements-pass' };
		}

		return { valid: true };
	},
};

export const ruleMaxLength = (max: number): ValidationRule => ({
	validate: (value: string) => {
		if (value.length > max) {
			return {
				valid: false,
				message: 'message.rule-value-max-length',
				argv: `${max}`,
			};
		}

		return { valid: true };
	},
});

export const VALIDATION_PRESENT = {
	UINT8: numberRange(0, 255),
	UINT16: numberRange(0, 65535),
	WORDS: ruleOnlyLetters,
	DIGITS: ruleOnlyDigits,
	DEVICE_PASSWORD: rulePassword,
	BEACON_TIME_SEC_MIN: numberRange(60, 2147483647),
	BEACON_TIME_SEC_MAX: numberRange(180, 2147483647),
	BEACON_SLEEP_TIME: numberRange(330, 2147483647),
	BEACON_SLEEP_TIME_ERROR: numberRange(300, 2147483647),
	AIN_TRUE_LOW: numberRange(0, 60000),
	DEVICE_NAME: ruleMaxLength(15),
};
