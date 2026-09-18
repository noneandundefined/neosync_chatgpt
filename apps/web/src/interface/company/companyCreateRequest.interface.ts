export interface CompanyCreateRequest {
	name: string;
	ttl: string;
	priority_model?: string;
	existing_task_action: string;
	launch_mode: string;
	incompatible_action: string;
	configuration_cource: 'device_sources' | 'template_sources' | 'file_sources';
	configuration_device_cource?: string | null;
	configuration_template_id?: number | null;
	configuration_file_name?: string | null;
	configuration_file_base64?: string | null;
}
