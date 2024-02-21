// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "description": "Log in to user and authenticate with the backend",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "User login information",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.authLoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.authLoginResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/auth/renew": {
            "post": {
                "description": "Renews token, preventing timeouts",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Renew token",
                "parameters": [
                    {
                        "description": "Token to renew",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.authRenewRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.authRenewResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/course/department/{dep}": {
            "get": {
                "description": "Queries for all UF courses in a three-letter department prefix",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "course"
                ],
                "summary": "Get courses in department",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Three-letter department prefix",
                        "name": "dep",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "400": {
                        "description": "Bad Request"
                    },
                    "503": {
                        "description": "Service Unavailable"
                    }
                }
            }
        },
        "/course/group/{dep}/{code}": {
            "get": {
                "description": "Gets (or creates) group of specified course",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "course"
                ],
                "summary": "Get group of specified course",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Three-letter department prefix",
                        "name": "dep",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Four-digit course code",
                        "name": "code",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    },
                    "503": {
                        "description": "Service Unavailable"
                    }
                }
            }
        },
        "/group/all": {
            "get": {
                "description": "Gets all public groups",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Gets groups",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.groupAllResponseItem"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/group/get/{id}": {
            "get": {
                "description": "Gets group information",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Get group",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Group ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.groupGetResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/user/bio": {
            "post": {
                "description": "Updates a user's bio",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Updates bio",
                "parameters": [
                    {
                        "description": "User token and new bio",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.userBioRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/user/get/{ident}": {
            "get": {
                "description": "Fetches user information and groups",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Fetch user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User identifier",
                        "name": "ident",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.userGetResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/user/join": {
            "post": {
                "description": "Adds a user to a group",
                "tags": [
                    "user"
                ],
                "summary": "Join group",
                "parameters": [
                    {
                        "description": "User token and group to join",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.userJoinRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/user/register": {
            "post": {
                "description": "Registers a new user given the provided arguments",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User registration information",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.userRegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.userRegisterResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "main.authLoginRequest": {
            "type": "object",
            "properties": {
                "ident": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "main.authLoginResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "main.authRenewRequest": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "main.authRenewResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "main.groupAllResponseItem": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "main.groupGetResponse": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "main.userBioRequest": {
            "type": "object",
            "properties": {
                "bio": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "main.userGetResponse": {
            "type": "object",
            "properties": {
                "bio": {
                    "type": "string"
                },
                "display_name": {
                    "type": "string"
                },
                "groups": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.userGetResponseGroup"
                    }
                },
                "ident": {
                    "type": "string"
                }
            }
        },
        "main.userGetResponseGroup": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "main.userJoinRequest": {
            "type": "object",
            "properties": {
                "group_id": {
                    "type": "integer"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "main.userRegisterRequest": {
            "type": "object",
            "properties": {
                "display_name": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "main.userRegisterResponse": {
            "type": "object",
            "properties": {
                "ident": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
