// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Shoyo",
            "url": "https://github.com/shoyo10/taskapi"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/tasks": {
            "get": {
                "description": "list all tasks",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/http.listTaskResp"
                        }
                    }
                }
            },
            "post": {
                "description": "create a task",
                "parameters": [
                    {
                        "description": "task fields",
                        "name": "reqBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.createTaskReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "task id",
                        "schema": {
                            "$ref": "#/definitions/http.createTaskResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errors.HTTPError"
                        }
                    }
                }
            }
        },
        "/tasks/{id}": {
            "put": {
                "description": "update a task",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "task id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "update task fields",
                        "name": "reqBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.updateTaskReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errors.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/errors.HTTPError"
                        }
                    }
                }
            },
            "delete": {
                "description": "delete a task",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "task id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errors.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/errors.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "errors.HTTPError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "message": {}
            }
        },
        "http.createTaskReq": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 1
                },
                "status": {
                    "type": "integer",
                    "enum": [
                        0,
                        1
                    ]
                }
            }
        },
        "http.createTaskResp": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/http.createTaskRespData"
                }
            }
        },
        "http.createTaskRespData": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "http.listTaskResp": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/http.listTaskRespData"
                    }
                }
            }
        },
        "http.listTaskRespData": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "status": {
                    "$ref": "#/definitions/model.EnumTaskStatus"
                }
            }
        },
        "http.updateTaskReq": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 1
                },
                "status": {
                    "type": "integer",
                    "enum": [
                        0,
                        1
                    ]
                }
            }
        },
        "model.EnumTaskStatus": {
            "type": "integer",
            "enum": [
                0,
                1
            ],
            "x-enum-varnames": [
                "TaskStatusIncomplete",
                "TaskStatusCompleted"
            ]
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:9090",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Task API Document",
	Description:      "This is task api document.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
