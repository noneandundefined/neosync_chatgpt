export type DocTab = 'tables' | 'functions';

export interface AnalyticTableDoc {
	titleKey: string;
	table: string;
}

export interface AnalyticFunctionParam {
	name: string;
	type: string;
	defaultValue?: string;
}

export interface AnalyticFunctionReturn {
	name: string;
	type: string;
}

export interface AnalyticFunctionDoc {
	name: string;
	signature: string;
	descriptionKey: string;
	parameters: AnalyticFunctionParam[];
	returns: AnalyticFunctionReturn[];
}

export const ANALYTIC_TABLES: AnalyticTableDoc[] = [
	{ titleKey: 'label.analytics-doc-table-users', table: 'analytics_users' },
	{ titleKey: 'label.analytics-doc-table-configurations-prof', table: 'configurations_prof' },
	{ titleKey: 'label.analytics-doc-table-records', table: 'analytic_records' },
	{ titleKey: 'label.analytics-doc-table-product-events', table: 'product_analytics_events' },
	{ titleKey: 'label.analytics-doc-table-configurations', table: 'configurations' },
	{ titleKey: 'label.analytics-doc-table-devices', table: 'devices' },
	{ titleKey: 'label.analytics-doc-table-syncs', table: 'syncs' },
	{ titleKey: 'label.analytics-doc-table-groups', table: 'groups' },
	{ titleKey: 'label.analytics-doc-table-group-members', table: 'group_members' },
	{ titleKey: 'label.analytics-doc-table-user-cores', table: 'user_cores' },
	{ titleKey: 'label.analytics-doc-table-user-roles', table: 'user_roles' },
	{ titleKey: 'label.analytics-doc-table-user-contacts', table: 'user_contacts' },
];

export const ANALYTIC_FUNCTIONS: AnalyticFunctionDoc[] = [
	{
		name: 'analytic_jsonb_int_bytes',
		signature: 'analytic_jsonb_int_bytes(p_data JSONB)',
		descriptionKey: 'label.analytics-doc-fn-jsonb-int-bytes',
		parameters: [{ name: 'p_data', type: 'JSONB' }],
		returns: [
			{ name: 'byte_val', type: 'INT' },
			{ name: 'pos', type: 'BIGINT' },
		],
	},
	{
		name: 'analytic_byte_slots',
		signature: 'analytic_byte_slots(p_data JSONB, p_bytes_per_slot INT DEFAULT 6, p_max_slots INT DEFAULT 8)',
		descriptionKey: 'label.analytics-doc-fn-byte-slots',
		parameters: [
			{ name: 'p_data', type: 'JSONB' },
			{ name: 'p_bytes_per_slot', type: 'INT', defaultValue: '6' },
			{ name: 'p_max_slots', type: 'INT', defaultValue: '8' },
		],
		returns: [
			{ name: 'slot', type: 'INT' },
			{ name: 'not_empty', type: 'BOOLEAN' },
		],
	},
	{
		name: 'analytic_typed_byte_slots',
		signature: 'analytic_typed_byte_slots(p_address_list JSONB, p_type_mask INT, p_bytes_per_slot INT DEFAULT 6, p_max_slots INT DEFAULT 8)',
		descriptionKey: 'label.analytics-doc-fn-typed-byte-slots',
		parameters: [
			{ name: 'p_address_list', type: 'JSONB' },
			{ name: 'p_type_mask', type: 'INT' },
			{ name: 'p_bytes_per_slot', type: 'INT', defaultValue: '6' },
			{ name: 'p_max_slots', type: 'INT', defaultValue: '8' },
		],
		returns: [
			{ name: 'slot', type: 'INT' },
			{ name: 'not_empty', type: 'BOOLEAN' },
			{ name: 'is_ble', type: 'BOOLEAN' },
			{ name: 'is_rs485', type: 'BOOLEAN' },
		],
	},
];
