include "console.iol"
include "interface.iol"
include "file.iol"

inputPort MyInput {
	Location: "socket://localhost:8000"
	Protocol: soap {
		.wsdl = "./myWsdl.wsdl";
		.wsdl.port = "MyInputServicePort";
		.dropRootValue = true
	}
	Interfaces: MyInterface
}

inputPort MyWSDL {
	Location: "socket://localhost:8001"
	Protocol: http {
        format = "raw"
        contentType = "application/xml"
        osc << {
            wsdl << {
                template = "/WSDL"
                method = "get"
            }
        }
    }
	Interfaces: HTTPInterface
}

main {
	while (true) {
		[
        		wsdl()(responseFile) {
        			file << {
        				.filename = "myWsdl.wsdl";
        			}
        			readFile@File( file )( responseFile )
        		}
        	]
		[sum(terms)(res) {
			res.result = terms.term1 + terms.term2
		}]
		[multiply(terms)(res) {
			res.result = terms.term1 * terms.term2
		}]
		[average(numbers)(res) {
			s = 0.0
			for(i = 0, i < #numbers.values, i++) {
				s += numbers.values[i]
			}
			res.result = s / #numbers.values
		}]
	}
}
