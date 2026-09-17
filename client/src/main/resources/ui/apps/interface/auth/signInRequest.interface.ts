export interface SignInRequest {
	login: string;
	password: string;
	remember_me: boolean;
}

export interface SignInRequestWithToken extends SignInRequest {
	turnstile_token: string;
}
