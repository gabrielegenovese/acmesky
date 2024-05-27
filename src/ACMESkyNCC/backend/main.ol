include "console.iol"
include "file.iol"
include "math.iol"
include "time.iol"

include "interface.iol"

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

init {
    global.nccs = {}
    global.bookings = {}
}

execution { concurrent }

main {
	[
		wsdl()(responseFile) {
			file << {
				.filename = "myWsdl.wsdl";
			}
			readFile@File(file)(responseFile)
		}
	]
	[add(request)(ncc) {
		id = new
		global.nccs.(id) << {
			.id = id,
			.name = request.name,
			.price = request.price,
			.location = request.location,
		}
		ncc << global.nccs.(id)
	}]
	[get()(nccs) {
		i = 0
		foreach (index : global.nccs) {
			nccs.nccs[i] << global.nccs.(index)
			i++
		}
	}]
	[getId(uuid)(ncc) {
		if (!is_defined(global.nccs.(uuid.value))) {
			println@Console("getId Error: NCC not found")()
			throw (NCCNotFound, {.message = uuid.value})
		}
		ncc << global.nccs.(uuid.value)
	}]
	[book(bookingRequest)(booking) {
		scope (book) {
			install (
				NCCNotFound => {
					println@Console("book Error: NCC not found")()
					booking << {.success = false}
				},
				OverlappingDate => {
					println@Console("book Error: Overlapping date")()
					booking << {.success = false}
				}
			);
			if (!is_defined(global.nccs.(bookingRequest.nccId))) {
				throw (NCCNotFound)
			} else {
				for (otherBooking in global.bookings.(bookingRequest.nccId)) {
					getTimestampFromString@Time(bookingRequest.date)(time1)
					getTimestampFromString@Time(otherBooking.date)(time2)
					abs@Math(int(time1 - time2))(timeDiff)
					if (timeDiff < 2 * 60 * 60 * 1000) { // Less than 2 hours
						throw (OverlappingDate)
					}
				}
				global.bookings.(bookingRequest.nccId)[#global.bookings.(bookingRequest.nccId)] << bookingRequest
				booking << {.success = true}
			}
		}
	}]
}
