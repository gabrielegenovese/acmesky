export interface Payment {
	id: string;
	user: string;
	amount: BigInteger;
	description: string;
	paid: boolean;
}
