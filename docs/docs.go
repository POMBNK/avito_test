// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Uchaev Roman",
            "url": "https://github.com/POMBNK",
            "email": "uchaevroman11@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/original_reports/{userID}": {
            "post": {
                "description": "Receiving CSV report with all user actions with segments by user ID, year and month.\nReport v1 exist according to the terms of reference about CSV reports",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "reports"
                ],
                "summary": "Get user's report v1",
                "operationId": "get-report-v1",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "userID",
                        "name": "userID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Month and year as the beginning of the time period for the CSV report",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/segment.ReportDateDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    }
                }
            }
        },
        "/api/reports/{userID}": {
            "post": {
                "description": "Receiving CSV report with all user actions with segments by user ID, year and month.\nReport v2 has better query and format to read.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "reports"
                ],
                "summary": "Get user's report v2",
                "operationId": "get-report-v2",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "userID",
                        "name": "userID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Month and year as the beginning of the time period for the CSV report",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/segment.ReportDateDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    }
                }
            }
        },
        "/api/segments": {
            "put": {
                "description": "Adds and removes tags to the user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "segments"
                ],
                "summary": "Edit segments to user",
                "operationId": "edit-segment",
                "parameters": [
                    {
                        "description": "User's segment info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/segment.SegmentsUsers"
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
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    }
                }
            },
            "post": {
                "description": "create segment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "segments"
                ],
                "summary": "Create segment",
                "operationId": "create-segment",
                "parameters": [
                    {
                        "description": "segment info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/segment.ToCreateSegmentDTO"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    }
                }
            },
            "delete": {
                "description": "delete segment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "segments"
                ],
                "summary": "Delete segment",
                "operationId": "delete-segment",
                "parameters": [
                    {
                        "description": "segment info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/segment.ToDeleteSegmentDTO"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    }
                }
            }
        },
        "/api/segments/ttl": {
            "post": {
                "description": "CronJobSegments is a function that handles cron job requests for checking segment TTL.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ttl"
                ],
                "summary": "Cron job to check ttl segments",
                "operationId": "ttl",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    }
                }
            }
        },
        "/api/segments/{userID}": {
            "get": {
                "description": "Get all active user's segments by userID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "segments"
                ],
                "summary": "Get active segments by id",
                "operationId": "get-active-segments",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "userID",
                        "name": "userID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/segment.ActiveSegments"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/apierror.ApiError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "apierror.ApiError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "developer_msg": {
                    "type": "string"
                }
            }
        },
        "segment.ActiveSegments": {
            "type": "object",
            "properties": {
                "ID": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "segment.ReportDateDTO": {
            "type": "object",
            "properties": {
                "month": {
                    "type": "string"
                },
                "year": {
                    "type": "string"
                }
            }
        },
        "segment.SegmentsUsers": {
            "type": "object",
            "properties": {
                "ID": {
                    "type": "string"
                },
                "add": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "name": {
                                "type": "string"
                            },
                            "ttl_days": {
                                "type": "integer"
                            }
                        }
                    }
                },
                "delete": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "userID": {
                    "type": "string"
                }
            }
        },
        "segment.ToCreateSegmentDTO": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "percent": {
                    "type": "integer"
                }
            }
        },
        "segment.ToDeleteSegmentDTO": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Avito Segment Service API",
	Description:      "API Server for Avito Segment Service",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
