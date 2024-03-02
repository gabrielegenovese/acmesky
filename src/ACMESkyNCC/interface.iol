type number: void {
	result: double
}

type pair: void {
    term1: double
    term2: double
}

type numbers: void {
    values[1,*]: double
}

interface MyInterface {
RequestResponse:
	sum(pair)(number),
	multiply(pair)(number),
	average(numbers)(number)
}

interface HTTPInterface {
RequestResponse:
	wsdl(void)(string)
}
