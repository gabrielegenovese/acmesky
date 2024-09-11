export interface Payment {
	id: string;
	created_at: string;
	updated_at: string;
	user: string;
	description: string;
	amount: number;
	link: string;
	paid: boolean;
}
