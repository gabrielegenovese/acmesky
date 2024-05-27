type nccRequest: void {
	name: string
	price: string
	location: string
}
type ncc: void {
	id: string
	name: string
	price: string
	location: string
}
type nccList: void {
    .nccs[0,*]: ncc
}
type uuid: void {
    .value: string
}
type error: void {
	.message: string
}
type bookingRequest: void {
	.nccId: string
	.name: string
	.date: string
}
type booking: void {
	.success: bool
}

interface MyInterface {
RequestResponse:
	add(nccRequest)(ncc),
	get(void)(nccList),
	getId(uuid)(ncc) throws NCCNotFound(error),
	book(bookingRequest)(booking)
}

interface HTTPInterface {
RequestResponse:
	wsdl(void)(string)
}
