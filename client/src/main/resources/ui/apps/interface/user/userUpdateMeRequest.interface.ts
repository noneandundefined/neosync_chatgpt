export interface UserUpdateMeRequest {
	email: string;
	phone: string;
	password?: string;
	turnstile_token: string;
}
