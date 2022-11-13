/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface EntitiesUser {
  /** @example "2022-06-05T14:26:02.302718+03:00" */
  created_at?: string
  /**
   * gorm:"uniqueIndex"
   * @example "name@email.com"
   */
  email?: string
  /** @example "WB7DRDWrJZRGbYrv2CKGkqbzvqdC" */
  id?: string
  /** @example "John Doe" */
  name?: string
  /** @example "2022-06-05T14:26:10.303278+03:00" */
  updated_at?: string
}

export interface RequestsCloudEvent {
  data?: any
  datacontenttype?: string
  id?: string
  source?: string
  specversion?: string
  time?: string
  type?: string
}

export interface ResponsesBadRequest {
  /** @example "The request body is not a valid JSON string" */
  data?: string
  /** @example "The request isn't properly formed" */
  message?: string
  /** @example "error" */
  status?: string
}

export interface ResponsesInternalServerError {
  /** @example "We ran into an internal error while handling the request." */
  message?: string
  /** @example "error" */
  status?: string
}

export interface ResponsesNoContent {
  /** @example "phone deleted successfully" */
  message?: string
  /** @example "success" */
  status?: string
}

export interface ResponsesOkEntitiesUser {
  data?: EntitiesUser
  /** @example "Request handled successfully" */
  message?: string
  /** @example "ok" */
  status?: string
}

export interface ResponsesUnauthorized {
  /** @example "Make sure your API key is set in the [X-API-Key] header in the request" */
  data?: string
  /** @example "You are not authorized to carry out this request." */
  message?: string
  /** @example "error" */
  status?: string
}

export interface ResponsesUnprocessableEntity {
  data?: Record<string, string[]>
  /** @example "validation errors while sending message" */
  message?: string
  /** @example "error" */
  status?: string
}
