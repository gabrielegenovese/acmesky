// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplatebank = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Gabriele Genovese",
            "email": "gabriele.genovese2@studio.unibo.it"
        },
        "license": {
            "name": "GPLv2",
            "url": "https://www.gnu.org/licenses/old-licenses/gpl-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/payment/new": {
            "put": {
                "description": "Create a new unpaid payment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payment"
                ],
                "summary": "New payment",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.Payment"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.Res"
                        }
                    }
                }
            }
        },
        "/payment/pay/{id}": {
            "post": {
                "description": "Pay an unpaid payment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payment"
                ],
                "summary": "Pay a payment",
                "parameters": [
                    {
                        "type": "string",
                        "description": "payment ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.Res"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.Res"
                        }
                    }
                }
            }
        },
        "/payment/{id}": {
            "get": {
                "description": "Given a payment ID, find the corresponding payment and return it.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payment"
                ],
                "summary": "Get a payment",
                "parameters": [
                    {
                        "type": "string",
                        "description": "payment ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.Res"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.Res"
                        }
                    }
                }
            },
            "delete": {
                "description": "Given a payment ID, find the corresponding payment and delete it.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payment"
                ],
                "summary": "Delete a payment",
                "parameters": [
                    {
                        "type": "string",
                        "description": "payment ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.Res"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.Res"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.Payment": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "$ref": "#/definitions/sql.NullTime"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "paid": {
                    "type": "boolean"
                },
                "updated_at": {
                    "type": "string"
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "api.Res": {
            "type": "object",
            "properties": {
                "res": {
                    "type": "string"
                }
            }
        },
        "sql.NullTime": {
            "type": "object",
            "properties": {
                "time": {
                    "type": "string"
                },
                "valid": {
                    "description": "Valid is true if Time is not NULL",
                    "type": "boolean"
                }
            }
        }
    }
}`

// SwaggerInfobank holds exported Swagger Info so clients can modify it
var SwaggerInfobank = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Bank Service API",
	Description:      "This is a minimal microservice to act as a bank.",
	InfoInstanceName: "bank",
	SwaggerTemplate:  docTemplatebank,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfobank.InstanceName(), SwaggerInfobank)
}